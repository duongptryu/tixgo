package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"tixgo/components"
	"tixgo/config"
	templatePort "tixgo/modules/template/ports"
	userPort "tixgo/modules/user/ports"

	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/duongptryu/gox/auth"
	"github.com/duongptryu/gox/database"
	"github.com/duongptryu/gox/logger"
	"github.com/duongptryu/gox/messaging"
	"github.com/duongptryu/gox/server/httpserver"
	"github.com/duongptryu/gox/syserr"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize logger first
	logger.Init(&logger.Config{
		Level:     slog.LevelInfo,
		Output:    os.Stdout,
		AddSource: false,
	})

	ctx := context.Background()
	logger.Info(ctx, "Starting TixGo API Server...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(ctx, "Failed to load configuration", logger.F("error", err))
	}

	logger.Info(ctx, "Configuration loaded successfully",
		logger.F("environment", cfg.App.Environment),
		logger.F("debug_mode", cfg.App.DebugMode))

	// Connect to database
	db, err := connectDatabase(ctx, &cfg.Database)
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to database", logger.F("error", err))
	}
	defer db.Close()

	logger.Info(ctx, "Database connected successfully")

	// Run migrations
	if err := runMigrations(ctx, db, &cfg.Database); err != nil {
		logger.Fatal(ctx, "Failed to run migrations", logger.F("error", err))
	}

	// Initialize app context
	appCtx, err := setupAppCtx(ctx, cfg, db)
	if err != nil {
		logger.Fatal(ctx, "Failed to initialize app context", logger.F("error", err))
	}

	// register event handlers
	registerEventHandlers(appCtx)

	// Setup HTTP server using server package
	srv := setupHTTPServer(ctx, cfg, appCtx)

	// Start server with graceful shutdown
	startServer(ctx, srv)
}

func connectDatabase(ctx context.Context, cfg *config.Database) (*sqlx.DB, error) {
	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	// Connect to database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func runMigrations(ctx context.Context, db *sqlx.DB, cfg *config.Database) error {
	logger.Info(ctx, "Running database migrations...")

	// Get SQL database instance for migrations
	sqlDB := db.DB

	// Create migration manager
	migrationManager, err := database.NewMigrationManager(sqlDB, &database.Config{
		Host:         cfg.Host,
		Port:         cfg.Port,
		User:         cfg.User,
		Password:     cfg.Password,
		Name:         cfg.Name,
		SSLMode:      cfg.SSLMode,
		Type:         cfg.Type,
		MaxOpenConns: cfg.MaxOpenConns,
		MaxIdleConns: cfg.MaxIdleConns,
		MaxLifetime:  cfg.MaxLifetime,
		MaxIdleTime:  cfg.MaxIdleTime,
	}, cfg.MigrationPath)
	if err != nil {
		return fmt.Errorf("failed to create migration manager: %w", err)
	}

	// Run migrations up
	if err := migrationManager.Up(); err != nil {
		// Check if it's "no change" error, which is acceptable
		if errors.Is(syserr.UnwrapError(err), migrate.ErrNoChange) {
			logger.Info(ctx, "No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info(ctx, "Database migrations completed successfully")
	return nil
}

func setupAppCtx(ctx context.Context, cfg *config.AppConfig, db *sqlx.DB) (components.AppContext, error) {
	jwtService := auth.NewJWTService(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
	)

	// init publisher
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kafkaSub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               cfg.Kafka.Brokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         "tixgo_consumer_group",
		},
		watermill.NewSlogLogger(logger.GetLogger()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka subscriber: %w", err)
	}

	kafkaPub, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   cfg.Kafka.Brokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewSlogLogger(logger.GetLogger()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka publisher: %w", err)
	}

	messagingBus, err := messaging.NewBus(messaging.Config{
		Publisher:  kafkaPub,
		Subscriber: kafkaSub,
		Logger:     logger.GetLogger(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create messaging bus: %w", err)
	}

	return components.NewAppContext(db, jwtService, messagingBus, messagingBus, messagingBus), nil
}

func setupHTTPServer(ctx context.Context, cfg *config.AppConfig, appCtx components.AppContext) *httpserver.Server {
	logger.Info(ctx, "Setting up HTTP server...")

	// Setup router with configuration
	router := httpserver.SetupRouter(httpserver.RouterConfig{
		ServiceName: cfg.App.Name,
		Environment: cfg.App.Environment,
		EnableCORS:  true,
		EnableAuth:  true,
	})

	// Register module routes
	registerRoutes(router, appCtx)

	// Create server with configuration
	srv := httpserver.New(httpserver.Config{
		Host:         cfg.Server.Host,
		Port:         cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}, router)

	logger.Info(ctx, "HTTP server configured",
		logger.F("address", srv.Addr()))

	return srv
}

func registerRoutes(router *gin.Engine, appCtx components.AppContext) {
	v1 := router.Group("/v1")
	// Register user module routes
	{
		userPort.RegisterUserRoutes(v1, appCtx)
		templatePort.RegisterTemplateRoutes(v1, appCtx)
	}

	// Add any additional module routes here
}

func registerEventHandlers(appCtx components.AppContext) {
	userPort.RegisterUserEventHandlers(appCtx.GetDispatcher(), appCtx)
}

func startServer(ctx context.Context, srv *httpserver.Server) {
	// Start server with graceful shutdown (blocks until shutdown)
	if err := srv.Start(ctx); err != nil {
		logger.Fatal(ctx, "Server failed", logger.F("error", err))
	}
}

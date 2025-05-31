package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"tixgo/config"
	"tixgo/modules/user"
	"tixgo/shared/database"
	"tixgo/shared/logger"
	"tixgo/shared/server/httpserver"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize logger first
	logger.Init(&logger.Config{
		Level:     slog.LevelInfo,
		Output:    os.Stdout,
		AddSource: true,
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

	// Initialize modules
	modules, err := initializeModules(ctx, db, cfg)
	if err != nil {
		logger.Fatal(ctx, "Failed to initialize modules", logger.F("error", err))
	}

	// Setup HTTP server using server package
	srv := setupHTTPServer(ctx, cfg, modules)

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
	migrationManager, err := database.NewMigrationManager(sqlDB, cfg)
	if err != nil {
		return fmt.Errorf("failed to create migration manager: %w", err)
	}

	// Run migrations up
	if err := migrationManager.Up(); err != nil {
		// Check if it's "no change" error, which is acceptable
		if err.Error() == "no change" {
			logger.Info(ctx, "No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info(ctx, "Database migrations completed successfully")
	return nil
}

type Modules struct {
	UserModule *user.Module
}

func initializeModules(ctx context.Context, db *sqlx.DB, cfg *config.AppConfig) (*Modules, error) {
	logger.Info(ctx, "Initializing application modules...")

	// JWT configuration - using hardcoded values for now
	// TODO: Move to config file
	jwtConfig := config.JWT{
		SecretKey:          cfg.JWT.SecretKey,
		AccessTokenExpiry:  cfg.JWT.AccessTokenExpiry,
		RefreshTokenExpiry: cfg.JWT.RefreshTokenExpiry,
	}

	// Initialize user module
	userModule, err := user.NewModule(db, jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user module: %w", err)
	}

	logger.Info(ctx, "User module initialized successfully")

	return &Modules{
		UserModule: userModule,
	}, nil
}

func setupHTTPServer(ctx context.Context, cfg *config.AppConfig, modules *Modules) *httpserver.Server {
	logger.Info(ctx, "Setting up HTTP server...")

	// Setup router with configuration
	router := httpserver.SetupRouter(httpserver.RouterConfig{
		Environment: cfg.App.Environment,
		EnableCORS:  true,
		EnableAuth:  true,
	})

	// Register module routes
	registerRoutes(router, modules)

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

func registerRoutes(router *gin.Engine, modules *Modules) {
	// Register user module routes
	modules.UserModule.RegisterRoutes(router)

	// Add any additional module routes here
}

func startServer(ctx context.Context, srv *httpserver.Server) {
	// Start server with graceful shutdown (blocks until shutdown)
	if err := srv.Start(ctx); err != nil {
		logger.Fatal(ctx, "Server failed", logger.F("error", err))
	}
}

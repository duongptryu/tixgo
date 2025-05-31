# TixGo API Server

This is the main API server for the TixGo application, implemented following the [Wild Workouts Go DDD example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example) patterns and architecture, specifically the [common server module](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/tree/master/internal/common/server).

## Architecture

The API server follows **Domain-Driven Design (DDD)** and **Clean Architecture** principles:

- **Domain Layer**: Business logic and entities
- **Application Layer**: Use cases and application services  
- **Infrastructure Layer**: Database repositories and external services
- **Interface Layer**: HTTP handlers and middleware
- **Shared Layer**: Common utilities including server management

## Features

### Core Components

- **Configuration Management**: Uses Viper for YAML configuration with environment variable overrides
- **Structured Logging**: JSON logging with context support using slog
- **Database**: PostgreSQL with sqlx and automatic migrations
- **HTTP Framework**: Gin with comprehensive middleware stack
- **Authentication**: JWT-based auth with access/refresh tokens
- **Server Management**: Centralized server utilities following Wild Workouts patterns
- **Graceful Shutdown**: Proper server shutdown handling with signal management

### Middleware Stack

1. **Request Context**: Adds request/operation IDs for traceability
2. **Request Logger**: Structured HTTP request logging
3. **Recovery**: Panic recovery with error logging
4. **CORS**: Cross-origin request support
5. **Error Handler**: Centralized error handling

### Modules

- **User Module**: Complete user management (registration, auth, profiles)
- **Extensible**: Easy to add new modules following the same patterns

## Quick Start

### Prerequisites

- Go 1.23.4+
- PostgreSQL database
- Make sure the database `tixgo_dev` exists

### Configuration

The server uses `config.yaml` for configuration:

```yaml
app:
  name: tixgo
  environment: dev
  debug_mode: true

server:
  host: localhost
  port: 8080
  read_timeout: 10s
  write_timeout: 10s
  idle_timeout: 10s

database: 
  type: postgres
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: tixgo_dev
  ssl_mode: disable
  max_open_conns: 10
  max_idle_conns: 5
  max_lifetime: 3600s
  max_idle_time: 3600s
  migration_path: file://migrations
```

### Building and Running

```bash
# Build the server
go build -o bin/api_server ./cmd/api_server

# Run the server
./bin/api_server
```

## Server Package Integration

The server now uses the `shared/server` package following Wild Workouts patterns:

### Server Setup

```go
// Setup router with configuration
routerConfig := server.RouterConfigFromAppConfig(cfg)
router := server.SetupRouter(routerConfig)

// Register module routes
registerRoutes(router, modules)

// Create server with configuration
serverConfig := server.ConfigFromAppConfig(cfg)
srv := server.New(serverConfig, router)

// Start server with graceful shutdown (blocks until shutdown)
if err := srv.Start(ctx); err != nil {
    logger.Fatal(ctx, "Server failed", logger.F("error", err))
}
```

### Benefits of Server Package

- **Centralized Configuration**: Type-safe server configuration conversion
- **Standardized Router Setup**: Consistent middleware pipeline across services
- **Automatic Health Endpoints**: `/health`, `/ready`, `/live` endpoints
- **Graceful Shutdown**: Signal handling and proper resource cleanup
- **Wild Workouts Compliance**: Follows exact patterns from the reference implementation

## API Endpoints

### Health Checks

- `GET /health` - Basic health check with timestamp
- `GET /ready` - Readiness check (service ready to handle requests)
- `GET /live` - Liveness check (service is alive)

### User Management

- `POST /api/v1/users/register` - User registration
- `POST /api/v1/users/verify-otp` - Email verification
- `POST /api/v1/users/login` - User login
- `GET /api/v1/users/profile` - Get user profile (requires auth)

## Wild Workouts Compliance

This implementation follows Wild Workouts patterns with enhanced server utilities:

### ✅ Configuration
- Viper-based configuration management
- Environment-specific config files
- Validation with struct tags
- **Type-safe server configuration conversion**

### ✅ Logging
- Structured JSON logging with slog
- Context-aware logging with request/operation IDs
- Centralized logger initialization

### ✅ Database
- sqlx for database operations
- Automatic migration management
- Connection pooling configuration

### ✅ HTTP Server
- **Centralized server utilities package**
- **Standardized router configuration**
- **Automatic health endpoint registration**
- Gin framework with middleware pipeline
- Graceful shutdown with context timeout

### ✅ Error Handling
- Custom error types with codes
- Centralized error handling middleware
- Structured error responses

### ✅ Architecture
- Clean modular structure
- Domain-driven design principles
- Hexagonal architecture (ports & adapters)
- **Common server utilities following Wild Workouts patterns**

## Development

### Adding New Modules

1. Create module structure in `modules/[module-name]/`
2. Implement domain, application, infrastructure, and interface layers
3. Add module initialization in `initializeModules()`
4. Register routes in `registerRoutes()`

### Using Server Package

```go
// Create router with standard configuration
routerConfig := server.RouterConfigFromAppConfig(appConfig)
router := server.SetupRouter(routerConfig)

// Add API groups
v1 := server.AddAPIGroup(router, "v1")
v1.POST("/login", loginHandler)

// Add protected routes
protected := server.AddProtectedGroup(v1, "/users")
protected.Use(middleware.RequireAuth(jwtService))
protected.GET("/profile", profileHandler)
```

### Environment Variables

Override config values using environment variables with `APP_` prefix:

```bash
export APP_DATABASE_HOST=localhost
export APP_DATABASE_PORT=5432
export APP_SERVER_PORT=8080
```

## Production Considerations

- [x] **Server utilities following Wild Workouts patterns**
- [x] **Centralized configuration management**
- [x] **Standardized middleware pipeline**
- [x] **Automatic health check endpoints**
- [x] **Graceful shutdown handling**
- [ ] Add proper JWT secret management
- [ ] Implement Redis for OTP storage
- [ ] Add rate limiting middleware
- [ ] Set up proper logging aggregation
- [ ] Configure SSL/TLS termination
- [ ] Add metrics and monitoring
- [ ] Implement database read replicas
- [ ] Add caching layers

## References

- [Wild Workouts Go DDD Example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
- [Wild Workouts Server Package](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/tree/master/internal/common/server)
- [Three Dots Labs Articles](https://threedots.tech/)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) 
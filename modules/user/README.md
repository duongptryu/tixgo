# User Module

A complete user management module implemented following **Wild Workouts Domain-Driven Design (DDD)** patterns with Clean Architecture and CQRS.

## Features

- ✅ User registration with email
- ✅ OTP-based email verification
- ✅ JWT-based authentication
- ✅ User profile management
- ✅ Password hashing with bcrypt
- ✅ Role-based user types (customer, organizer, admin)
- ✅ Comprehensive error handling
- ✅ Clean Architecture with DDD patterns

## Architecture

This module follows the **Wild Workouts DDD example** architecture with clear separation of concerns:

```
modules/user/
├── domain/           # Domain layer - pure business logic
│   ├── user.go      # User aggregate root
│   ├── repository.go # Repository interfaces (ports)
│   ├── errors.go    # Domain-specific errors
│   └── events.go    # Domain events
├── app/             # Application layer - CQRS
│   ├── command/     # Commands (write operations)
│   │   ├── register_user.go
│   │   ├── verify_otp.go
│   │   └── login_user.go
│   ├── query/       # Queries (read operations)
│   │   └── get_user_profile.go
│   └── service.go   # Application service facade
├── adapters/        # Infrastructure layer (adapters)
│   ├── user_postgres.go  # PostgreSQL repository
│   └── otp_store.go     # In-memory OTP store
├── ports/           # Interface layer (ports)
│   └── http.go      # HTTP/REST API handlers
└── module.go        # Dependency injection & wiring
```

### Layer Responsibilities

1. **Domain Layer**: Pure business logic, no external dependencies
2. **Application Layer**: Use cases, commands, queries (CQRS)
3. **Infrastructure Layer**: External concerns (database, storage)
4. **Interface Layer**: External communication (HTTP, gRPC)

## API Endpoints

### 1. Register User
```http
POST /api/v1/users/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "user_type": "customer"
}
```

**Response:**
```json
{
  "data": {
    "user_id": 1,
    "message": "User registered successfully. Please verify your email with the OTP.",
    "otp": "123456"
  }
}
```

### 2. Verify OTP
```http
POST /api/v1/users/verify-otp
Content-Type: application/json

{
  "email": "user@example.com",
  "otp": "123456"
}
```

**Response:**
```json
{
  "data": {
    "user_id": 1,
    "message": "Email verified successfully."
  }
}
```

### 3. Login User
```http
POST /api/v1/users/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "data": {
    "user_id": 1,
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

### 4. Get User Profile
```http
GET /api/v1/users/profile
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "user_type": "customer",
    "status": "active",
    "email_verified": true,
    "created_at": "2024-01-01T00:00:00Z",
    "last_login": "2024-01-01T12:00:00Z"
  }
}
```

## Usage

### Basic Setup

```go
package main

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
    usermodule "tixgo/modules/user"
)

func main() {
    // Database connection
    db, err := sqlx.Connect("postgres", "your-database-url")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Configure user module
    userConfig := usermodule.Config{
        JWTSecretKey:          "your-secret-key",
        JWTAccessTokenExpiry:  15 * time.Minute,
        JWTRefreshTokenExpiry: 24 * time.Hour,
    }

    // Initialize user module
    userModule, err := usermodule.NewModule(db, userConfig)
    if err != nil {
        log.Fatal(err)
    }

    // Setup router and register routes
    router := gin.Default()
    userModule.RegisterRoutes(router)

    // Start server
    router.Run(":8080")
}
```

### Using the Application Service Directly

```go
// Register a user
result, err := userModule.UserService.RegisterUser(ctx, command.RegisterUserCommand{
    Email:     "user@example.com",
    Password:  "password123",
    FirstName: "John",
    LastName:  "Doe",
    UserType:  "customer",
})

// Login a user
loginResult, err := userModule.UserService.LoginUser(ctx, command.LoginUserCommand{
    Email:    "user@example.com",
    Password: "password123",
})

// Get user profile
profile, err := userModule.UserService.GetUserProfile(ctx, query.GetUserProfileQuery{
    UserID: 1,
})
```

## Configuration

The module requires the following configuration:

```go
type Config struct {
    JWTSecretKey          string        // Secret key for JWT signing
    JWTAccessTokenExpiry  time.Duration // Access token expiry (e.g., 15 minutes)
    JWTRefreshTokenExpiry time.Duration // Refresh token expiry (e.g., 24 hours)
}
```

## Database Schema

The module uses the existing `users` table from the database migrations:

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    date_of_birth DATE,
    user_type user_type_enum DEFAULT 'customer',
    status user_status_enum DEFAULT 'active',
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);
```

## Error Handling

The module uses the existing `shared/syserr` package for structured error handling:

- `invalid_argument` - Invalid input data
- `validation_error` - Validation failures
- `unauthorized` - Authentication failures
- `forbidden` - Authorization failures
- `not_found` - Resource not found
- `conflict` - Resource conflicts (e.g., email already exists)
- `internal` - Internal server errors

## Security Features

- **Password Hashing**: bcrypt with default cost (12)
- **JWT Tokens**: Access and refresh token pairs
- **Email Verification**: OTP-based email verification
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: Parameterized queries with sqlx

## Testing

Run the tests:

```bash
# Test domain layer
go test ./modules/user/domain/...

# Test application layer
go test ./modules/user/app/...

# Test infrastructure layer
go test ./modules/user/adapters/...

# Test all
go test ./modules/user/...
```

## Production Considerations

1. **OTP Delivery**: Replace in-memory OTP store with Redis and implement email/SMS delivery
2. **Database**: Use connection pooling and read replicas for scalability
3. **JWT Security**: Use RS256 instead of HS256 for better security
4. **Rate Limiting**: Implement rate limiting for authentication endpoints
5. **Monitoring**: Add metrics and logging for observability
6. **Caching**: Cache user profiles for better performance

## Dependencies

The module leverages existing shared utilities:

- `shared/auth` - JWT token management
- `shared/syserr` - Structured error handling
- `shared/response` - HTTP response formatting
- `shared/context` - Request context utilities

## Wild Workouts Alignment

This implementation follows the [Wild Workouts Go DDD example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example) patterns:

- ✅ Domain-Driven Design with aggregates
- ✅ Clean Architecture with dependency inversion
- ✅ CQRS pattern for read/write separation
- ✅ Hexagonal architecture (ports & adapters)
- ✅ Repository pattern for data access
- ✅ Domain events for integration
- ✅ Dependency injection for testability 
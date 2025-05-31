# User Module Documentation

## Overview

The User Module is implemented following the **Wild Workouts Domain-Driven Design (DDD)** architecture pattern, providing a clean separation of concerns and a modular structure that can easily be extracted into a microservice.

## Architecture Layers

### 1. Domain Layer (`internal/modules/user/domain/`)

**Pure business logic with no external dependencies**

- **`user.go`**: User aggregate root with business logic
- **`events.go`**: Domain events for event-driven architecture
- **`repository.go`**: Repository interface (port)
- **`errors.go`**: Domain-specific errors

#### Key Features:
- Rich domain model with encapsulated business logic
- Password hashing and validation
- User type management (customer, organizer, admin)
- Email verification workflow
- Domain events for integration

### 2. Application Layer (`internal/modules/user/app/`)

**CQRS (Command Query Responsibility Segregation) pattern**

#### Commands (`app/command/`)
- **`register_user.go`**: User registration
- **`login_user.go`**: User authentication with JWT
- **`update_profile.go`**: Profile updates

#### Queries (`app/query/`)
- **`get_user_profile.go`**: Retrieve user profile
- **`list_users.go`**: List users with pagination and filtering

#### Application Service (`app/service.go`)
- Orchestrates commands and queries
- Provides facade methods for easier access

### 3. Infrastructure Layer (`internal/modules/user/adapters/`)

**External dependencies and data access**

- **`user_postgres.go`**: PostgreSQL repository implementation
- Implements domain repository interfaces
- Handles data mapping between domain and database

### 4. Interface Layer (`internal/modules/user/ports/`)

**External communication interfaces**

- **`http.go`**: REST API endpoints with Gin framework
- **`grpc.go`**: gRPC service implementation
- Both share the same application layer

### 5. Module Bootstrap (`internal/modules/user/module.go`)

**Dependency injection and module wiring**

- Wires all dependencies together
- Provides module configuration
- Exposes necessary services to other modules

## Features Implemented

### REST API Endpoints

#### Public Endpoints
- `POST /api/v1/users/register` - User registration
- `POST /api/v1/users/login` - User authentication

#### Protected Endpoints (JWT required)
- `GET /api/v1/users/profile` - Get current user profile
- `PUT /api/v1/users/profile` - Update user profile

### gRPC Services

- `ListUsers` - List users with pagination and filtering
- `GetUserByID` - Retrieve user by ID

### Authentication & Authorization

- **JWT-based authentication** with access and refresh tokens
- **Password hashing** using bcrypt with cost 12
- **Role-based access control** (customer, organizer, admin)
- **Email verification** workflow support

### Domain Events

The module publishes domain events for integration:

- `UserRegistered`
- `UserPasswordChanged`
- `UserProfileUpdated`
- `UserEmailVerified`
- `UserSuspended`
- `UserActivated`
- `UserLoggedIn`

## Usage Example

### 1. Module Initialization

```go
package main

import (
    "time"
    usermodule "tixgo/internal/modules/user"
)

func main() {
    // Configure user module
    userConfig := usermodule.Config{
        JWTSecretKey:          "your-secret-key-here",
        JWTAccessTokenExpiry:  15 * time.Minute,
        JWTRefreshTokenExpiry: 24 * time.Hour,
    }

    // Initialize user module
    userModule, err := usermodule.NewModule(db, userConfig)
    if err != nil {
        log.Fatal("Failed to initialize user module:", err)
    }

    // Register HTTP routes
    userModule.HTTPHandler.RegisterRoutes(router)

    // Register gRPC service
    userpb.RegisterUserServiceServer(grpcServer, userModule.GRPCServer)
}
```

### 2. REST API Examples

#### Register User
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "user_type": "customer"
  }'
```

#### Login User
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### Get Profile
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer your_jwt_token_here"
```

### 3. gRPC Examples

#### List Users
```bash
grpcurl -plaintext \
  -d '{
    "page": 1,
    "page_size": 10,
    "user_type": "USER_TYPE_CUSTOMER",
    "search": "john"
  }' \
  localhost:9090 user.v1.UserService/ListUsers
```

#### Get User by ID
```bash
grpcurl -plaintext \
  -d '{"id": "user-uuid-here"}' \
  localhost:9090 user.v1.UserService/GetUserByID
```

## Database Schema

The user module expects the following PostgreSQL table structure:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('customer', 'organizer', 'admin')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login TIMESTAMP NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_user_type ON users(user_type);
CREATE INDEX idx_users_status ON users(status);
```

## Configuration

### Environment Variables

```env
# Database
DATABASE_URL=postgresql://username:password@localhost/tixgo?sslmode=disable

# JWT
JWT_SECRET_KEY=your-very-secure-secret-key-here
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=24h

# Server
HTTP_PORT=8080
GRPC_PORT=9090
```

## Testing

### Unit Tests
Each layer has comprehensive unit tests:

```bash
# Test domain layer
go test ./internal/modules/user/domain/...

# Test application layer
go test ./internal/modules/user/app/...

# Test infrastructure layer
go test ./internal/modules/user/adapters/...

# Test interfaces
go test ./internal/modules/user/ports/...
```

### Integration Tests
```bash
# Test the entire module
go test ./internal/modules/user/...
```

## Best Practices Implemented

1. **Domain-Driven Design**: Clear separation between domain, application, and infrastructure
2. **CQRS Pattern**: Separate command and query responsibilities
3. **Repository Pattern**: Abstracted data access layer
4. **Event-Driven Architecture**: Domain events for loose coupling
5. **Clean Architecture**: Dependencies point inward
6. **Interface Segregation**: Small, focused interfaces
7. **Dependency Injection**: Configurable and testable dependencies
8. **Error Handling**: Structured error handling with proper HTTP/gRPC status codes
9. **Security**: Password hashing, JWT tokens, input validation
10. **Documentation**: Comprehensive API documentation

## Migration to Microservice

When ready to extract this module into a separate microservice:

1. **Extract the module directory** to a new repository
2. **Add main.go** (example provided in `cmd/example/main.go`)
3. **Update imports** to reflect the new module path
4. **Add configuration management** (environment variables, config files)
5. **Add health checks and monitoring**
6. **Update database connection** to use dedicated database
7. **Add API gateway** for external communication
8. **Implement service discovery** if using microservice orchestration

## Future Enhancements

- [ ] Email verification service integration
- [ ] Password reset functionality
- [ ] OAuth2/OIDC integration
- [ ] Rate limiting
- [ ] User activity logging
- [ ] Account lockout mechanisms
- [ ] GDPR compliance features (data export, deletion)
- [ ] Multi-factor authentication
- [ ] Social login integration
- [ ] User preferences and settings 
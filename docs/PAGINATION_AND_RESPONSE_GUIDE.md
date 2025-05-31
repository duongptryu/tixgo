# Simplified Pagination and Response Guide

## Overview

This guide covers the simplified pagination and response system in TixGo. The system provides clean, simple APIs with minimal boilerplate code - no complex metadata, timestamps, or heavy abstractions.

## Architecture

```
internal/shared/
├── pagination/
│   ├── types.go      # Simple Paging struct with basic methods
│   ├── http.go       # HTTP query parsing helpers  
│   └── grpc.go       # gRPC conversion helpers
└── response/
    ├── types.go      # Simple success and error response types
    ├── success.go    # Success response helpers
    └── error.go      # Error response helpers
```

## Features

### ✅ **Simple Pagination**
- Single `Paging` struct with `page`, `limit`, `total`, `cursor`, `next_cursor`
- Automatic defaults: page=1, limit=10
- Easy offset calculation: `GetOffset()` method
- Basic pagination info: `HasNext()`, `HasPrev()`, `GetTotalPages()`

### ✅ **Clean Responses**
- Simple JSON structure: `{data, paging, filter}`
- No complex metadata or timestamps
- Automatic syserr error mapping
- Multiple response formats supported

## Quick Start

### 1. **HTTP REST API Usage**

```go
// In your HTTP handler
func (h *Handler) ListUsers(c *gin.Context) {
    // Parse pagination from query parameters
    paging := pagination.ParsePagingFromQuery(c)
    
    // Parse filters
    filter := UserFilter{
        Search: c.Query("search"),
        Status: c.Query("status"),
    }

    // Your business logic
    users, total, err := h.userService.ListUsers(c.Request.Context(), paging, filter)
    if err != nil {
        response.ErrorFromSyserrJSON(c, err)
        return
    }

    // Set total count
    paging.Total = total

    // Return paginated response
    response.SuccessPaginatedJSON(c, users, paging)
}
```

### 2. **gRPC Usage**

```go
// In your gRPC server
func (s *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    // Parse pagination from gRPC request
    paging := pagination.ParsePagingFromGRPC(req.Page, req.Limit)

    // Your business logic
    users, total, err := s.userService.ListUsers(ctx, paging, filter)
    if err != nil {
        return nil, err
    }

    // Convert pagination response
    page, limit, _ := paging.ToPagingGRPC()

    return &pb.ListUsersResponse{
        Users: users,
        Page:  page,
        Limit: limit,
        Total: total,
    }, nil
}
```

## Core Types

### **Paging Struct**

```go
type Paging struct {
    Page       int   `json:"page" form:"page"`
    Limit      int   `json:"limit" form:"limit"`
    Total      int64 `json:"total" form:"total"`
    FakeCursor int   `json:"cursor" form:"cursor"`
    NextCursor int   `json:"next_cursor"`
}

// Methods
paging.Fulfill()              // Apply defaults (page=1, limit=10)
paging.GetOffset()            // Calculate database offset
paging.HasNext()              // Check if has next page
paging.HasPrev()              // Check if has previous page
paging.GetTotalPages()        // Calculate total pages
```

### **Response Types**

```go
// Simple success response
type successRes struct {
    Data   interface{} `json:"data"`
    Paging interface{} `json:"paging,omitempty"`
    Filter interface{} `json:"filter,omitempty"`
}

// Simple error response
type errorRes struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}
```

## Usage Examples

### **HTTP Integration**

```go
// Parse pagination from query
paging := pagination.ParsePagingFromQuery(c) // ?page=2&limit=20

// Database usage
offset := paging.GetOffset() // (page-1) * limit
query := "SELECT * FROM users LIMIT $1 OFFSET $2"
db.Query(query, paging.Limit, offset)

// Response helpers
response.SuccessJSON(c, data)                           // Simple success
response.SuccessPaginatedJSON(c, data, paging)          // With pagination
response.SuccessWithFilterJSON(c, data, filter)        // With filter
response.SuccessWithPagingAndFilterJSON(c, data, paging, filter) // Both
```

### **gRPC Integration**

```go
// Parse from gRPC request
paging := pagination.ParsePagingFromGRPC(req.Page, req.Limit)

// Convert back to gRPC
page, limit, total := paging.ToPagingGRPC()
```

### **Service Layer Pattern**

```go
type UserService interface {
    ListUsers(ctx context.Context, paging pagination.Paging, filter UserFilter) ([]*User, int64, error)
}

func (s *UserServiceImpl) ListUsers(ctx context.Context, paging pagination.Paging, filter UserFilter) ([]*User, int64, error) {
    return s.userRepo.List(ctx, paging, filter)
}
```

### **Repository Layer Pattern**

```go
func (r *Repository) List(ctx context.Context, paging pagination.Paging, filter UserFilter) ([]*User, int64, error) {
    // Use pagination
    offset := paging.GetOffset()
    limit := paging.Limit
    
    // Data query
    query := "SELECT * FROM users WHERE name ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3"
    rows, err := r.db.Query(query, "%"+filter.Search+"%", limit, offset)
    
    // Count query
    countQuery := "SELECT COUNT(*) FROM users WHERE name ILIKE $1"
    var total int64
    r.db.QueryRow(countQuery, "%"+filter.Search+"%").Scan(&total)
    
    return users, total, nil
}
```

## Response Formats

### **Simple Success Response**
```json
{
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### **Paginated Response**
```json
{
  "data": [
    {"id": 1, "name": "John Doe", "email": "john@example.com"},
    {"id": 2, "name": "Jane Smith", "email": "jane@example.com"}
  ],
  "paging": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "cursor": 0,
    "next_cursor": 0
  }
}
```

### **Response with Filter**
```json
{
  "data": [...],
  "paging": {...},
  "filter": {
    "search": "john",
    "status": "active"
  }
}
```

### **Error Response**
```json
{
  "code": "validation_error",
  "message": "Invalid input",
  "details": {
    "email": "Valid email is required"
  }
}
```

## Available Response Functions

### **Success Responses**
```go
response.SuccessJSON(c, data)
response.SuccessCreatedJSON(c, data)
response.SuccessPaginatedJSON(c, data, paging)
response.SuccessWithFilterJSON(c, data, filter)
response.SuccessWithPagingAndFilterJSON(c, data, paging, filter)
response.SuccessNoContentJSON(c)
```

### **Error Responses**
```go
response.ErrorJSON(c, statusCode, code, message)
response.ErrorWithDetailsJSON(c, statusCode, code, message, details)
response.ErrorFromSyserrJSON(c, err)  // Automatic mapping
response.ValidationErrorJSON(c, message, fieldErrors)

// Shortcuts
response.BadRequestJSON(c, message)
response.UnauthorizedJSON(c, message)
response.ForbiddenJSON(c, message)
response.NotFoundJSON(c, message)
response.ConflictJSON(c, message)
response.InternalServerErrorJSON(c, message)
```

## Migration from Complex System

### **Before (Complex)**
```go
// Old complex approach
config := pagination.DefaultConfig()
paginationReq, err := pagination.ParseFromQuery(c, config)
if err != nil {
    return err
}

result, err := service.List(ctx, ListQuery{
    Pagination: paginationReq,
    Search: search,
})

resp := response.Success(result.Items)
resp.WithPagination(&result.Pagination)
resp.WithRequestID(requestID)
resp.JSON(c, http.StatusOK)
```

### **After (Simple)**
```go
// New simple approach
paging := pagination.ParsePagingFromQuery(c)

users, total, err := service.ListUsers(ctx, paging, filter)
if err != nil {
    response.ErrorFromSyserrJSON(c, err)
    return
}

paging.Total = total
response.SuccessPaginatedJSON(c, users, paging)
```

## Best Practices

### ✅ **DO:**

1. **Use simple parsing**
   ```go
   paging := pagination.ParsePagingFromQuery(c)
   ```

2. **Set total count after service call**
   ```go
   users, total, err := service.List(ctx, paging, filter)
   paging.Total = total
   ```

3. **Use appropriate response helpers**
   ```go
   response.SuccessPaginatedJSON(c, data, paging)
   ```

4. **Handle errors with automatic mapping**
   ```go
   response.ErrorFromSyserrJSON(c, err)
   ```

### ❌ **DON'T:**

1. **Manual offset calculations**
   ```go
   offset := (page - 1) * limit // Use paging.GetOffset() instead
   ```

2. **Complex metadata handling**
   ```go
   // Avoid - no need for timestamps, request IDs, etc.
   ```

3. **Manual status code mapping**
   ```go
   // Use response.ErrorFromSyserrJSON(c, err) instead
   ```

## Testing

### **Unit Tests**
```go
func TestPaging(t *testing.T) {
    paging := pagination.Paging{Page: 2, Limit: 10}
    paging.Fulfill()
    
    assert.Equal(t, 2, paging.Page)
    assert.Equal(t, 10, paging.Limit)
    assert.Equal(t, 10, paging.GetOffset())
}
```

### **Integration Tests**
```go
func TestListAPI(t *testing.T) {
    resp := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/v1/users?page=2&limit=5", nil)
    
    router.ServeHTTP(resp, req)
    
    assert.Equal(t, 200, resp.Code)
    
    var result struct {
        Data   []User            `json:"data"`
        Paging pagination.Paging `json:"paging"`
    }
    
    json.Unmarshal(resp.Body.Bytes(), &result)
    assert.Equal(t, 2, result.Paging.Page)
    assert.Equal(t, 5, result.Paging.Limit)
}
```

## Complete Example

```go
// Handler
func (h *UserHandler) ListUsers(c *gin.Context) {
    paging := pagination.ParsePagingFromQuery(c)
    filter := UserFilter{Search: c.Query("search")}
    
    users, total, err := h.service.ListUsers(c.Request.Context(), paging, filter)
    if err != nil {
        response.ErrorFromSyserrJSON(c, err)
        return
    }
    
    paging.Total = total
    response.SuccessPaginatedJSON(c, users, paging)
}

// Service
func (s *UserService) ListUsers(ctx context.Context, paging pagination.Paging, filter UserFilter) ([]*User, int64, error) {
    return s.repo.List(ctx, paging, filter)
}

// Repository
func (r *UserRepo) List(ctx context.Context, paging pagination.Paging, filter UserFilter) ([]*User, int64, error) {
    offset := paging.GetOffset()
    
    query := "SELECT * FROM users WHERE name ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3"
    rows, err := r.db.Query(query, "%"+filter.Search+"%", paging.Limit, offset)
    
    // ... scan rows
    
    var total int64
    r.db.QueryRow("SELECT COUNT(*) FROM users WHERE name ILIKE $1", "%"+filter.Search+"%").Scan(&total)
    
    return users, total, nil
}
```

This simplified approach provides all the necessary functionality while being much easier to understand and maintain. 
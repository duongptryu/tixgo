# Pagination Guide

## Overview

The TixGo project provides a comprehensive pagination system in `internal/common/pagination/` that can be reused across all modules. This follows the DDD architecture and ensures consistency across REST API, gRPC, and database operations.

## Architecture

```
internal/common/pagination/
├── types.go    # Core pagination types and logic
├── http.go     # HTTP/REST specific helpers  
└── grpc.go     # gRPC specific helpers
```

## Core Types

### **Request**
Represents pagination parameters from clients:
```go
type Request struct {
    Page     int `json:"page" form:"page" binding:"min=1"`
    PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}
```

### **Response** 
Represents pagination metadata in responses:
```go
type Response struct {
    Page       int  `json:"page"`
    PageSize   int  `json:"page_size"`
    TotalCount int  `json:"total_count"`
    TotalPages int  `json:"total_pages"`
    HasNext    bool `json:"has_next"`
    HasPrev    bool `json:"has_previous"`
}
```

### **Filter**
Represents database-level pagination:
```go
type Filter struct {
    Page     int
    PageSize int
    Offset   int
}
```

## Usage Examples

### **1. HTTP/REST API**

#### In HTTP Handler
```go
func (h *Handler) ListItems(c *gin.Context) {
    // Parse pagination from query parameters
    paginationReq, err := pagination.ParseFromQuery(c, pagination.DefaultConfig())
    if err != nil {
        h.handleError(c, err)
        return
    }

    // Create query with pagination
    query := query.ListItems{
        Pagination: paginationReq,
        // ... other filters
    }

    // Execute query
    result, err := h.app.ListItems(c.Request.Context(), query)
    if err != nil {
        h.handleError(c, err)
        return
    }

    // Create paginated response
    response := pagination.NewListResponse(result.Items, paginationReq, result.Pagination.TotalCount)

    // Add pagination headers (optional)
    result.Pagination.AddToHeader(c)

    c.JSON(http.StatusOK, response)
}
```

#### Query Parameters
```bash
GET /api/v1/items?page=2&page_size=20&search=golang
```

#### Response Format
```json
{
  "data": [...],
  "pagination": {
    "page": 2,
    "page_size": 20,
    "total_count": 150,
    "total_pages": 8,
    "has_next": true,
    "has_previous": true
  }
}
```

### **2. gRPC API**

#### In gRPC Server
```go
func (s *Server) ListItems(ctx context.Context, req *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
    // Parse pagination with defaults
    paginationReq, err := pagination.NewGRPCRequest(req.Page, req.PageSize, pagination.DefaultConfig())
    if err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }

    // Create query
    query := query.ListItems{
        Pagination: paginationReq,
        // ... other filters
    }

    // Execute query
    result, err := s.app.ListItems(ctx, query)
    if err != nil {
        return nil, s.handleError(err)
    }

    // Convert pagination response
    page, pageSize, totalCount := result.Pagination.ToGRPCResponse()

    return &pb.ListItemsResponse{
        Items:      convertItems(result.Items),
        Page:       page,
        PageSize:   pageSize,
        TotalCount: totalCount,
    }, nil
}
```

### **3. Application Layer (Query)**

#### Query Definition
```go
type ListItems struct {
    Pagination pagination.Request
    Category   *string
    Status     *ItemStatus
    Search     string
}

type ListItemsResult struct {
    Items      []*domain.Item
    Pagination pagination.Response
}
```

#### Query Handler
```go
func (h ListItemsHandler) Handle(ctx context.Context, query ListItems) (*ListItemsResult, error) {
    // Convert to repository filter
    filter := domain.ListItemsFilter{
        Page:     query.Pagination.Page,
        PageSize: query.Pagination.PageSize,
        Category: query.Category,
        Status:   query.Status,
        Search:   query.Search,
    }

    // Execute repository query
    items, totalCount, err := h.itemRepo.List(ctx, filter)
    if err != nil {
        return nil, err
    }

    return &ListItemsResult{
        Items:      items,
        Pagination: pagination.NewResponse(query.Pagination, totalCount),
    }, nil
}
```

### **4. Repository Layer**

#### Repository Implementation
```go
func (r *PostgresItemRepository) List(ctx context.Context, filter domain.ListItemsFilter) ([]*domain.Item, int, error) {
    // Build WHERE conditions
    whereClause, args := r.buildWhereClause(filter)
    
    // Count total records
    countQuery := fmt.Sprintf("SELECT COUNT(*) FROM items %s", whereClause)
    var totalCount int
    err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
    if err != nil {
        return nil, 0, err
    }

    // Get paginated results
    offset := pagination.CalculateOffset(filter.Page, filter.PageSize)
    dataQuery := fmt.Sprintf(`
        SELECT * FROM items %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, len(args)+1, len(args)+2)
    
    args = append(args, filter.PageSize, offset)
    
    // Execute query and return results
    // ...
}
```

## Configuration

### **Default Configuration**
```go
config := pagination.DefaultConfig()
// config.DefaultPageSize = 20
// config.MaxPageSize = 100
```

### **Custom Configuration**
```go
config := pagination.Config{
    DefaultPageSize: 10,
    MaxPageSize:     50,
}
```

## Validation Rules

- **Page**: Must be ≥ 1 (defaults to 1)
- **PageSize**: Must be ≥ 1 and ≤ 100 (defaults to 20)
- **Automatic defaults**: Invalid values are automatically corrected

## Best Practices

### **1. Consistent Response Format**
Always use `pagination.NewListResponse()` for HTTP APIs:
```go
response := pagination.NewListResponse(data, paginationReq, totalCount)
```

### **2. Add Pagination Headers**
For better client experience:
```go
result.Pagination.AddToHeader(c)
```

### **3. Validate Early**
Validate pagination parameters at the interface layer:
```go
paginationReq, err := pagination.ParseFromQuery(c, pagination.DefaultConfig())
if err != nil {
    return err
}
```

### **4. Use Type-Safe Helpers**
For gRPC conversions:
```go
page, pageSize, totalCount := result.Pagination.ToGRPCResponse()
```

## Example: Complete Events Module

Here's how you'd implement pagination in a new Events module:

### **1. Query Definition**
```go
// internal/modules/events/app/query/list_events.go
type ListEvents struct {
    Pagination pagination.Request
    Category   *string
    Status     *domain.EventStatus
    DateFrom   *time.Time
    DateTo     *time.Time
    Search     string
}
```

### **2. HTTP Handler**
```go
// internal/modules/events/ports/http.go
func (h *HTTPHandler) ListEvents(c *gin.Context) {
    // Parse pagination
    paginationReq, err := pagination.ParseFromQuery(c, pagination.DefaultConfig())
    if err != nil {
        h.handleError(c, err)
        return
    }

    // Parse filters
    category := c.Query("category")
    search := c.Query("search")
    
    // Create and execute query
    query := query.ListEvents{
        Pagination: paginationReq,
        Category:   &category,
        Search:     search,
    }
    
    result, err := h.app.ListEvents(c.Request.Context(), query)
    if err != nil {
        h.handleError(c, err)
        return
    }

    // Return paginated response
    events := make([]EventResponse, len(result.Events))
    for i, event := range result.Events {
        events[i] = h.eventToResponse(event)
    }
    
    response := pagination.NewListResponse(events, paginationReq, result.Pagination.TotalCount)
    result.Pagination.AddToHeader(c)
    
    c.JSON(http.StatusOK, response)
}
```

### **3. gRPC Handler**
```go
// internal/modules/events/ports/grpc.go
func (s *GRPCServer) ListEvents(ctx context.Context, req *eventpb.ListEventsRequest) (*eventpb.ListEventsResponse, error) {
    paginationReq, err := pagination.NewGRPCRequest(req.Page, req.PageSize, pagination.DefaultConfig())
    if err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }

    query := query.ListEvents{
        Pagination: paginationReq,
        Search:     req.Search,
    }

    result, err := s.app.ListEvents(ctx, query)
    if err != nil {
        return nil, s.handleError(err)
    }

    page, pageSize, totalCount := result.Pagination.ToGRPCResponse()
    
    return &eventpb.ListEventsResponse{
        Events:     convertEvents(result.Events),
        Page:       page,
        PageSize:   pageSize,
        TotalCount: totalCount,
    }, nil
}
```

## Migration from Existing Code

To migrate existing pagination to use the common system:

### **Before**
```go
type ListUsers struct {
    Page     int
    PageSize int
    // other fields...
}
```

### **After**
```go
type ListUsers struct {
    Pagination pagination.Request
    // other fields...
}
```

The migration process:
1. Replace individual `Page`/`PageSize` fields with `pagination.Request`
2. Update handlers to use pagination helpers
3. Update responses to use `pagination.Response`
4. Use type-safe conversion methods for gRPC

## Future Enhancements

- [ ] Cursor-based pagination for large datasets
- [ ] Configurable pagination per module
- [ ] Pagination analytics and monitoring
- [ ] Custom sorting integration
- [ ] Performance optimization for count queries 
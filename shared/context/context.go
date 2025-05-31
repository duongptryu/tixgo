package context

import (
	"context"
	"tixgo/shared/auth"
)

// Context key types to avoid collisions
type contextKey string

const (
	// OperationIDKey is used for storing operation IDs in context
	OperationIDKey contextKey = "operationID"
	// RequestIDKey is used for storing request IDs in context
	RequestIDKey contextKey = "requestID"
	// UserIDKey is used for storing user IDs in context
	UserIDKey contextKey = "userID"
	// UserTypeKey is used for storing user types in context
	UserTypeKey contextKey = "userType"
	// AuthClaimsKey is used for storing auth claims in context
	AuthClaimsKey contextKey = "authClaims"
)

// Operation ID context utilities

// WithOperationID adds an operation ID to the context
func WithOperationID(ctx context.Context, operationID string) context.Context {
	if operationID == "" {
		return ctx
	}
	return context.WithValue(ctx, OperationIDKey, operationID)
}

// GetOperationID retrieves the operation ID from context
func GetOperationID(ctx context.Context) string {
	if value := ctx.Value(OperationIDKey); value != nil {
		if operationID, ok := value.(string); ok {
			return operationID
		}
	}
	return ""
}

// Request ID context utilities

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		return ctx
	}
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if value := ctx.Value(RequestIDKey); value != nil {
		if requestID, ok := value.(string); ok {
			return requestID
		}
	}
	return ""
}

// User ID context utilities

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	if userID == "" {
		return ctx
	}
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID retrieves the user ID from context
func GetUserIDFromContext(ctx context.Context) string {
	if value := ctx.Value(UserIDKey); value != nil {
		if userID, ok := value.(string); ok {
			return userID
		}
	}
	return ""
}

// User type context utilities

// WithUserType adds a user type to the context
func WithUserType(ctx context.Context, userType string) context.Context {
	if userType == "" {
		return ctx
	}
	return context.WithValue(ctx, UserTypeKey, userType)
}

// GetUserType retrieves the user type from context
func GetUserTypeFromContext(ctx context.Context) string {
	if value := ctx.Value(UserTypeKey); value != nil {
		if userType, ok := value.(string); ok {
			return userType
		}
	}
	return ""
}

func WithAuthClaims(ctx context.Context, claims *auth.Claims) context.Context {
	return context.WithValue(ctx, AuthClaimsKey, claims)
}

func GetAuthClaimsFromContext(ctx context.Context) *auth.Claims {
	if value := ctx.Value(AuthClaimsKey); value != nil {
		if claims, ok := value.(*auth.Claims); ok {
			return claims
		}
	}
	return nil
}

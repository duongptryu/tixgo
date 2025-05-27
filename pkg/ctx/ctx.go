package ctx

import (
	"context"
)

type operationIDKey string

const (
	OperationIDKey operationIDKey = "operationID"
)

func WithOperationID(ctx context.Context, operationID string) context.Context {
	if operationID == "" {
		return ctx
	}

	return context.WithValue(ctx, OperationIDKey, operationID)
}

func GetOperationID(ctx context.Context) string {
	value := ctx.Value(OperationIDKey)
	if value != nil {
		if operationID, ok := value.(string); ok {
			return operationID
		}
	}

	return "" // fallback to empty string
}

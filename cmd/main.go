package main

import (
	"context"

	"github.com/pkg/errors"

	pkgCtx "tixgo/pkg/ctx"
	"tixgo/pkg/logger"
)

func main() {
	ctx := context.Background()
	ctx = pkgCtx.WithOperationID(ctx, "123")

	err := errors.New("test error")
	logger.LogError(ctx, err)
}

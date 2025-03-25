package logger

import (
	"context"

	"go.uber.org/zap"
)

type loggerKeyCtx struct{}

func InsertLoggerToContext(ctx context.Context, logger *zap.Logger) context.Context {
	ctx = context.WithValue(ctx, loggerKeyCtx{}, logger)
	return ctx
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(loggerKeyCtx{}).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}
	return l
}

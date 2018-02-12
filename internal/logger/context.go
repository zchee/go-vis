package logger

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var contextKey struct{}

func NewContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

func WithContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	return NewContext(ctx, FromContext(ctx).With(fields...))
}

func FromContext(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(contextKey).(*zap.Logger)
	if !ok {
		return grpc_zap.Extract(ctx)
	}
	return l
}

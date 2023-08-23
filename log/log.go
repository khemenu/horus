package log

import (
	"context"
	"log/slog"
)

type logCtxKey struct{}

func FromCtx(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value(logCtxKey{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}

	return l
}

func WithCtx(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, logCtxKey{}, logger)
}

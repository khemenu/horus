package log

import (
	"context"
	"io"
	"log/slog"
)

type logCtxKey struct{}

// TODO: do nothing handler
var discard = slog.New(slog.NewTextHandler(io.Discard, nil))

func FromCtx(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value(logCtxKey{}).(*slog.Logger)
	if !ok {
		return discard
	}

	return l
}

func WithCtx(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, logCtxKey{}, logger)
}

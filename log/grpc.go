package log

import (
	"context"
	"log/slog"
	"time"

	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor(l *slog.Logger, level slog.Level) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		t := time.Now()
		l := l.With(slog.String("x", xid.New().String()))
		l.LogAttrs(ctx, level, "IN_", slog.String("method", info.FullMethod))

		ctx = Into(ctx, l)
		res, err := handler(ctx, req)
		if s, _ := status.FromError(err); err != nil {
			l = l.With(slog.Group("err",
				slog.Int("code", int(s.Code())),
				slog.String("name", s.Code().String()),
				slog.String("desc", s.Message()),
			))
		}

		l.LogAttrs(ctx, level, "OUT", slog.Duration("dt", time.Since(t)))
		return res, err
	}
}

package log

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/xid"
)

type resWriterWithStatusCapture struct {
	http.ResponseWriter
	code int
}

func (w *resWriterWithStatusCapture) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func HttpLogger(l *slog.Logger, level slog.Level, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		t := time.Now()

		ctx := req.Context()
		ctx = Into(ctx, l)

		l := l.With(slog.String("x", xid.New().String()))
		l.LogAttrs(ctx, level, "IN_", slog.String("method", req.Method), slog.String("uri", req.RequestURI))

		req = req.WithContext(ctx)
		res_ := &resWriterWithStatusCapture{ResponseWriter: res, code: 0}
		next.ServeHTTP(res_, req)

		l.LogAttrs(ctx, level, "OUT", slog.Int("code", res_.code), slog.Duration("dt", time.Since(t)))
	})
}

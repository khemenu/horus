package main

import (
	"log/slog"
	"net/http"
	"time"

	"khepri.dev/horus/log"
)

type responseWriter struct {
	http.ResponseWriter
	status_code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status_code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func httpLog(l *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()

		res := &responseWriter{w, http.StatusOK}
		w = res
		r = r.WithContext(log.WithCtx(r.Context(), l))

		l.Info("REST IN ", slog.String("method", r.Method), "url", r.URL)
		defer func() {
			l.Info("REST OUT", slog.Duration("dt", time.Since(t0)), slog.Int("status", res.status_code))
		}()

		next.ServeHTTP(w, r)
	})
}

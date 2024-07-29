package rabbitmq

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/pasztorpisti/qs"
	"khepri.dev/horus/log"
)

func ParseHttpBody[T any](w http.ResponseWriter, r *http.Request, v *T) error {
	l := log.From(r.Context())

	body, err := io.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		l.Error("read body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("read body: %w", err)
	}
	if err := qs.Unmarshal(v, string(body)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}

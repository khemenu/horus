package app

import (
	"errors"
	"fmt"
	"net/http"

	"khepri.dev/horus"
	"khepri.dev/horus/cmd/horus/server/frame"
)

type app struct {
	horus.Stores

	conf *horus.Config
}

func NewHorus(stores horus.Stores, conf *horus.Config) (horus.Horus, error) {
	if conf == nil {
		conf = &horus.Config{}
	}
	if err := conf.Normalize(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &app{
		Stores: stores,
		conf:   conf,
	}, nil
}

func (h *app) Config() *horus.Config {
	return h.conf
}

func (h *app) Verify(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := frame.FromCtx(r.Context()); ok {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(horus.CookieNameAccessToken)
		if err != nil {
			http.Error(w, "no access token", http.StatusUnprocessableEntity)
			return
		}

		access_token, err := h.Tokens().GetByValue(r.Context(), cookie.Value, horus.AccessToken)
		if err != nil {
			if errors.Is(err, horus.ErrNotExist) {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			http.Error(w, "failed to get token details", http.StatusInternalServerError)
			return
		}

		f := frame.NewFrame(h, access_token)
		ctx := frame.WithCtx(r.Context(), f)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

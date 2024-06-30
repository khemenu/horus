package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/log"
)

func HandleAuth(mux *http.ServeMux, svr horus.Server) {
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /auth/bearer", func(w http.ResponseWriter, r *http.Request) {
		l := log.From(r.Context())

		h := r.Header.Get("Authorization")
		if h == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Authorization header is missing")
			return
		}
		if !strings.HasPrefix(h, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Invalid authorization header format")
			return
		}

		v := h[len("Bearer "):]
		_, err := svr.Auth().TokenSignIn(r.Context(), &horus.TokenSignInRequest{
			Token: v,
		})
		if err != nil {
			s, ok := status.FromError(err)
			if ok {
				switch s.Code() {
				case codes.NotFound:
					fallthrough
				case codes.Unauthenticated:
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			l.Error("bearer ", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /auth/basic/sign-in", func(w http.ResponseWriter, r *http.Request) {
		l := log.From(r.Context())

		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		res, err := svr.Auth().BasicSignIn(r.Context(), &horus.BasicSignInRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			s, ok := status.FromError(err)
			if ok {
				switch s.Code() {
				case codes.NotFound:
					fallthrough
				case codes.Unauthenticated:
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			l.Error("basic sign-in", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  horus.TokenKeyName,
			Value: res.Token.Value,

			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),

			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
	})
	mux.HandleFunc("GET /auth/sign-out", func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(horus.TokenKeyName); err != nil {
			// No token found
		} else if _, err := svr.Auth().SignOut(r.Context(), &horus.SingOutRequest{Token: cookie.Value}); err == nil {
			// Ok
		} else {
			// TODO: log
		}

		http.SetCookie(w, &http.Cookie{
			Name: horus.TokenKeyName,

			Expires: time.Unix(0, 0),
			MaxAge:  -1,

			Path:   "/",
			Secure: true,
		})
	})
}

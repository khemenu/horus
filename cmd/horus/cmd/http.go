package cmd

import (
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/log"
)

func HandleAuth(mux *http.ServeMux, svc horus.Service) {
	mux.HandleFunc("/auth/basic/sign-in", func(w http.ResponseWriter, r *http.Request) {
		l := log.From(r.Context())

		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		res, err := svc.Auth().BasicSignIn(r.Context(), &horus.BasicSignInRequest{
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
	mux.HandleFunc("/auth/sign-out", func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(horus.TokenKeyName); err != nil {
			// No token found
		} else if _, err := svc.Auth().SignOut(r.Context(), &horus.SingOutRequest{Token: &horus.Token{Value: cookie.Value}}); err == nil {
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

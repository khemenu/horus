package cmd

import (
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus/ent/proto/khepri/horus"
	"khepri.dev/horus/service"
	"khepri.dev/horus/tokens"
)

func HandleAuth(mux *http.ServeMux, svc service.Service) {
	mux.HandleFunc("/auth/basic/sign-in", func(w http.ResponseWriter, r *http.Request) {
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

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  tokens.CookieName,
			Value: res.Token.Value,

			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),

			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
	})
	mux.HandleFunc("/auth/sign-out", func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(tokens.CookieName); err != nil {
			// No cookie found => no tokens
		} else if _, err := svc.Auth().SignOut(r.Context(), &horus.SingOutRequest{Token: &horus.Token{Value: cookie.Value}}); err == nil {
			// Ok
		} else {
			// TODO: log
		}

		http.SetCookie(w, &http.Cookie{
			Name: tokens.CookieName,

			Expires: time.Unix(0, 0),
			MaxAge:  -1,

			Path:   "/",
			Secure: true,
		})
	})
}

package cmd

import (
	"fmt"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/horus/cmd/rabbitmq"
	"khepri.dev/horus/log"
)

func HandleRabbitMqHttpAuth(mux *http.ServeMux, svr horus.Server) {
	mux.HandleFunc("POST /auth/rabbitmq/user", func(w http.ResponseWriter, r *http.Request) {
		l := log.From(r.Context()).With("handler", "rabbitmq")
		r = r.WithContext(log.Into(r.Context(), l))

		req := rabbitmq.UserReq{}
		if err := rabbitmq.ParseHttpBody(w, r, &req); err != nil {
			return
		}

		_, err := svr.Auth().BasicSignIn(r.Context(), &horus.BasicSignInRequest{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			st, _ := status.FromError(err)
			fmt.Printf("st: %v\n", st)
			switch st.Code() {
			case codes.Unauthenticated:
				fallthrough
			case codes.FailedPrecondition:
				rabbitmq.Deny(w)
				return

			default:
				l.Error("unexpected response from the auth service", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		rabbitmq.Allow(w)
	})
	mux.HandleFunc("POST /auth/rabbitmq/vhost", func(w http.ResponseWriter, r *http.Request) {
		l := log.From(r.Context()).With("handler", "rabbitmq")
		r = r.WithContext(log.Into(r.Context(), l))

		req := rabbitmq.VhostReq{}
		if err := rabbitmq.ParseHttpBody(w, r, &req); err != nil {
			return
		}

		// TODO: authz?
		// Allow any vhost for any user for now
		rabbitmq.Allow(w)
	})
}

package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	authnV1 "k8s.io/api/authentication/v1"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/log"
	"khepri.dev/horus/server/frame"
)

func HandleKubeWebhook(mux *http.ServeMux, svr horus.Server) {
	mux.HandleFunc("POST /auth/kube", func(w http.ResponseWriter, r *http.Request) {
		l := log.From(r.Context()).With("handler", "kube-webhook")

		review := authnV1.TokenReview{}
		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		supported_versions := []string{
			"authentication.k8s.io/v1",
			"authentication.k8s.io/v1beta1",
		}
		if !slices.Contains(supported_versions, review.APIVersion) {
			res := fmt.Sprintf(`given API version "%s" not supported.
Supported API versions are: %s
  - `,
				review.APIVersion,
				strings.Join(supported_versions, "\n  -"),
			)

			w.WriteHeader(http.StatusPreconditionFailed)
			w.Write([]byte(res))
			return
		}

		reviewed := authnV1.TokenReview{
			TypeMeta: review.TypeMeta,
			Status: authnV1.TokenReviewStatus{
				Authenticated: true,
			},
		}

		sign_in, err := svr.Auth().TokenSignIn(r.Context(), &horus.TokenSignInRequest{
			Token: review.Spec.Token,
		})
		if err != nil {
			reviewed.Status.Authenticated = false

			s, _ := status.FromError(err)
			switch s.Code() {
			case codes.InvalidArgument:
				fallthrough
			case codes.Unauthenticated:
				break

			default:
				l.Error("unexpected response from the auth service", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if !reviewed.Status.Authenticated {
			if err := json.NewEncoder(w).Encode(&reviewed); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		ctx := frame.WithContext(r.Context(), &frame.Frame{
			Actor: &ent.User{ID: uuid.UUID(sign_in.Token.Owner.Id)},
		})
		user, err := svr.User().Get(ctx, &horus.GetUserRequest{})
		if err != nil {
			l.Error("get user details", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		reviewed.Status.User = authnV1.UserInfo{
			Username: user.Alias,
			UID:      uuid.UUID(user.Id).String(),
		}
		if err := json.NewEncoder(w).Encode(&reviewed); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

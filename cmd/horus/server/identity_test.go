package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/pb"
)

func TestListIdentities(t *testing.T) {
	WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		identity, err := h.Identities().New(ctx, &horus.IdentityInit{
			OwnerId:    h.user.Id,
			Kind:       horus.IdentityMail,
			Value:      "second@khepri.dev",
			VerifiedBy: "khepri",
		})
		require.NoError(err)

		res, err := h.client.ListIdentities(ctx, &pb.ListIdentitiesReq{})
		require.NoError(err)
		require.ElementsMatch(
			[]string{string(h.identity.Value), string(identity.Value)},
			fx.MapV(res.Identities, func(v *pb.Identity) string { return v.Value }),
		)
	})(t)
}

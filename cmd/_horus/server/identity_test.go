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

func TestDeleteIdentity(t *testing.T) {
	t.Run("not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		_, err := h.client.DeleteIdentity(ctx, &pb.DeleteIdentityReq{IdentityValue: "not exist"})
		require.NoError(err)
	}))

	t.Run("exists", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		_, err := h.client.DeleteIdentity(ctx, &pb.DeleteIdentityReq{IdentityValue: string(h.identity.Value)})
		require.NoError(err)

		res, err := h.client.ListIdentities(ctx, &pb.ListIdentitiesReq{})
		require.NoError(err)
		require.Empty(res.Identities)
	}))

	t.Run("delete others one", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		identity, err := h.Identities().New(ctx, &horus.IdentityInit{
			Kind:       horus.IdentityMail,
			Value:      "other@khepri.dev",
			VerifiedBy: "khepri",
		})
		require.NoError(err)

		_, err = h.client.DeleteIdentity(ctx, &pb.DeleteIdentityReq{IdentityValue: string(identity.Value)})
		require.NoError(err)
	}))
}

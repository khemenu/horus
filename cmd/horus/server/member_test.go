package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/pb"
)

func TestAddMemberIdentity(t *testing.T) {
	t.Run("mine", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{
			OwnerId: h.user.Id,
		})
		require.NoError(err)

		_, err = h.client.AddMemberIdentity(ctx, &pb.AddMemberIdentityReq{
			MemberId:      rst.Owner.Id[:],
			IdentityValue: string(h.identity.Value),
		})
		require.NoError(err)
	}))

	t.Run("others", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{
			OwnerId: other.Id,
		})
		require.NoError(err)

		_, err = h.client.AddMemberIdentity(ctx, &pb.AddMemberIdentityReq{
			MemberId:      rst.Owner.Id[:],
			IdentityValue: string(h.identity.Value),
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.PermissionDenied, s.Code())
	}))

	t.Run("identity does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{
			OwnerId: h.user.Id,
		})
		require.NoError(err)

		_, err = h.client.AddMemberIdentity(ctx, &pb.AddMemberIdentityReq{
			MemberId:      rst.Owner.Id[:],
			IdentityValue: "not exist",
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code(), err)
	}))
}

func TestRemoveMemberIdentity(t *testing.T) {
	t.Run("mine", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{
			OwnerId: h.user.Id,
		})
		require.NoError(err)

		err = h.Members().AddIdentity(ctx, rst.Owner.Id, h.identity.Value)
		require.NoError(err)

		_, err = h.client.RemoveMemberIdentity(ctx, &pb.RemoveMemberIdentityReq{
			MemberId:      rst.Owner.Id[:],
			IdentityValue: string(h.identity.Value),
		})
		require.NoError(err)
	}))

	t.Run("others", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		identity, err := h.Identities().New(ctx, &horus.IdentityInit{
			Kind:       horus.IdentityMail,
			Value:      "amun@khepri.dev",
			VerifiedBy: "khepri",
		})
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{
			OwnerId: identity.OwnerId,
		})
		require.NoError(err)

		err = h.Members().AddIdentity(ctx, rst.Owner.Id, identity.Value)
		require.NoError(err)

		_, err = h.client.RemoveMemberIdentity(ctx, &pb.RemoveMemberIdentityReq{
			MemberId:      rst.Owner.Id[:],
			IdentityValue: string(h.identity.Value),
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.PermissionDenied, s.Code())
	}))
}

package server_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/pb"
)

func TestUpdateMember(t *testing.T) {
	t.Run("org does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		org_id := uuid.New()
		_, err := h.client.UpdateMember(ctx, &pb.UpdateMemberReq{
			OrgId: org_id[:],
			Name:  "foo",
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("not a member", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		_, err = h.client.UpdateMember(ctx, &pb.UpdateMemberReq{
			OrgId: rst.Org.Id[:],
			Name:  "foo",
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("as a member", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		member, err := h.Members().New(ctx, horus.MemberInit{
			OrgId:  rst.Org.Id,
			UserId: h.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)
		require.NotEqual("foo", member.Name)

		_, err = h.client.UpdateMember(ctx, &pb.UpdateMemberReq{
			OrgId: rst.Org.Id[:],
			Name:  "foo",
		})
		require.NoError(err)

		member, err = h.Members().GetById(ctx, member.Id)
		require.NoError(err)
		require.Equal("foo", member.Name)
	}))
}

func TestAddMemberIdentity(t *testing.T) {
	t.Run("org does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		org_id := uuid.New()
		_, err := h.client.AddMemberIdentity(ctx, &pb.AddMemberIdentityReq{
			OrgId:         org_id[:],
			IdentityValue: string(h.identity.Value),
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("identity does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		_, err = h.client.AddMemberIdentity(ctx, &pb.AddMemberIdentityReq{
			OrgId:         rst.Org.Id[:],
			IdentityValue: "not exist",
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("others identity", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		other_identity, err := h.Identities().New(ctx, &horus.IdentityInit{
			OwnerId:    other.Id,
			Kind:       horus.IdentityMail,
			Value:      "other@khepri.dev",
			VerifiedBy: "khepri",
		})
		require.NoError(err)
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		_, err = h.client.AddMemberIdentity(ctx, &pb.AddMemberIdentityReq{
			OrgId:         rst.Org.Id[:],
			IdentityValue: string(other_identity.Value),
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("my identity", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		_, err = h.client.AddMemberIdentity(ctx, &pb.AddMemberIdentityReq{
			OrgId:         rst.Org.Id[:],
			IdentityValue: string(h.identity.Value),
		})
		require.NoError(err)
	}))

	t.Run("identity already exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		_, err = h.client.AddMemberIdentity(ctx, &pb.AddMemberIdentityReq{
			OrgId:         rst.Org.Id[:],
			IdentityValue: string(h.identity.Value),
		})
		require.NoError(err)

		_, err = h.client.AddMemberIdentity(ctx, &pb.AddMemberIdentityReq{
			OrgId:         rst.Org.Id[:],
			IdentityValue: string(h.identity.Value),
		})
		require.NoError(err)
	}))
}

func TestRemoveMemberIdentity(t *testing.T) {
	t.Run("org does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		org_id := uuid.New()
		_, err := h.client.RemoveMemberIdentity(ctx, &pb.RemoveMemberIdentityReq{
			OrgId:         org_id[:],
			IdentityValue: string(h.identity.Value),
		})
		require.NoError(err)
	}))

	t.Run("identity does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		_, err = h.client.RemoveMemberIdentity(ctx, &pb.RemoveMemberIdentityReq{
			OrgId:         rst.Org.Id[:],
			IdentityValue: "not exist",
		})
		require.NoError(err)
	}))
}

package server_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/pb"
)

func TestListOrg(t *testing.T) {
	t.Run("org does not exist", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		res, err := h.client.ListOrgs(ctx, &pb.ListOrgsReq{})
		require.NoError(err)
		require.Empty(res.Orgs)
	}))

	t.Run("org exists", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		org1, err := h.client.NewOrg(ctx, &pb.NewOrgReq{})
		require.NoError(err)

		org2, err := h.client.NewOrg(ctx, &pb.NewOrgReq{})
		require.NoError(err)

		res, err := h.client.ListOrgs(ctx, &pb.ListOrgsReq{})
		require.NoError(err)
		require.ElementsMatch(
			[][]byte{org1.Org.Id, org2.Org.Id},
			fx.MapV(res.Orgs, func(v *pb.Org) []byte {
				return v.Id
			}))
	}))
}

func TestUpdateOrg(t *testing.T) {
	t.Run("invalid org ID", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		_, err := h.client.UpdateOrg(ctx, &pb.UpdateOrgReq{Org: &pb.Org{Id: []byte{42}}})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.InvalidArgument, s.Code())
	}))

	t.Run("org does not exist", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		id := uuid.New()
		_, err := h.client.UpdateOrg(ctx, &pb.UpdateOrgReq{Org: &pb.Org{Id: id[:]}})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("as an owner", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		org, err := h.client.NewOrg(ctx, &pb.NewOrgReq{})
		require.NoError(err)

		_, err = h.client.UpdateOrg(ctx, &pb.UpdateOrgReq{Org: &pb.Org{
			Id:   org.Org.Id,
			Name: "foo",
		}})
		require.NoError(err)
	}))

	t.Run("as a member", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		owner, err := h.Users().New(ctx)
		require.NoError(err)

		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: owner.Id})
		require.NoError(err)

		h.Members().New(ctx, horus.MemberInit{
			OrgId:  org.Id,
			UserId: h.user.Id,
			Role:   horus.RoleOrgMember,
		})

		_, err = h.client.UpdateOrg(ctx, &pb.UpdateOrgReq{Org: &pb.Org{
			Id:   org.Id[:],
			Name: "foo",
		}})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.PermissionDenied, s.Code())
	}))
}

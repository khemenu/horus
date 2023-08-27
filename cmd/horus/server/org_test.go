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

func TestInviteUser(t *testing.T) {
	t.Run("invalid argument", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		testCases := []struct {
			desc string
			req  *pb.InviteUserReq
		}{
			{
				desc: "invalid ID",
				req: &pb.InviteUserReq{
					OrgId: []byte{42},
					Identity: &pb.Identity{
						Kind:  pb.IdentityKind_IDENTITY_KIND_UNSPECIFIED,
						Value: "foo@khepri.dev",
					},
				},
			},
			{
				desc: "kind unspecified",
				req: &pb.InviteUserReq{
					OrgId: org.Id[:],
					Identity: &pb.Identity{
						Kind:  pb.IdentityKind_IDENTITY_KIND_UNSPECIFIED,
						Value: "foo@khepri.dev",
					},
				},
			},
			{
				desc: "invalid mail addres",
				req: &pb.InviteUserReq{
					OrgId: org.Id[:],
					Identity: &pb.Identity{
						Kind:  pb.IdentityKind_IDENTITY_KIND_MAIL,
						Value: "royale with cheese",
					},
				},
			},
		}
		for _, tC := range testCases {
			t.Log(tC.desc)
			_, err = h.client.InviteUser(ctx, tC.req)
			s, ok := status.FromError(err)
			require.True(ok)
			require.Equal(codes.InvalidArgument, s.Code())
		}
	}))

	t.Run("only owner can invite a user", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		owner, err := h.Users().New(ctx)
		require.NoError(err)

		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: owner.Id})
		require.NoError(err)

		req := &pb.InviteUserReq{
			OrgId: org.Id[:],
			Identity: &pb.Identity{
				Kind:  pb.IdentityKind_IDENTITY_KIND_MAIL,
				Value: "foo@khepri",
			},
		}

		{
			// Invite as non-member.
			_, err = h.client.InviteUser(ctx, req)
			s, ok := status.FromError(err)
			require.True(ok)
			require.Equal(codes.PermissionDenied, s.Code())
		}

		{
			// Invite as a member
			h.Members().New(ctx, horus.MemberInit{
				OrgId:  org.Id,
				UserId: h.user.Id,
				Role:   horus.RoleOrgMember,
			})

			_, err = h.client.InviteUser(ctx, req)
			s, ok := status.FromError(err)
			require.True(ok)
			require.Equal(codes.PermissionDenied, s.Code())
		}
	}))

	t.Run("user that does not exist", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		_, err = h.client.InviteUser(ctx, &pb.InviteUserReq{
			OrgId: org.Id[:],
			Identity: &pb.Identity{
				Kind:  pb.IdentityKind_IDENTITY_KIND_MAIL,
				Value: "foo@khepri.dev",
			},
		})
		require.NoError(err)

		identity, err := h.Identities().GetByValue(ctx, "foo@khepri.dev")
		require.NoError(err)
		require.Equal(horus.IdentityMail, identity.Kind)
		require.Equal(horus.Unverified, identity.VerifiedBy)
	}))

	t.Run("user already a member", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		_, err = h.client.InviteUser(ctx, &pb.InviteUserReq{
			OrgId: org.Id[:],
			Identity: &pb.Identity{
				Kind:  pb.IdentityKind_IDENTITY_KIND_MAIL,
				Value: h.identity.Value,
			},
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.AlreadyExists, s.Code())
	}))
}

func TestJoinOrg(t *testing.T) {
	t.Run("invalid argument", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		_, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		testCases := []struct {
			desc string
			req  *pb.JoinOrgReq
		}{
			{
				desc: "invalid ID",
				req: &pb.JoinOrgReq{
					OrgId: []byte{42},
				},
			},
		}
		for _, tC := range testCases {
			t.Log(tC.desc)
			_, err = h.client.JoinOrg(ctx, tC.req)
			s, ok := status.FromError(err)
			require.True(ok)
			require.Equal(codes.InvalidArgument, s.Code())
		}
	}))

	t.Run("already a member", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		_, err = h.client.JoinOrg(ctx, &pb.JoinOrgReq{OrgId: org.Id[:]})
		require.NoError(err)
	}))

	t.Run("not invited", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		_, err = h.client.JoinOrg(ctx, &pb.JoinOrgReq{OrgId: org.Id[:]})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.PermissionDenied, s.Code())
	}))

	t.Run("invited", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other := h.WithNewIdentity(ctx, &horus.IdentityInit{
			Kind:  horus.IdentityMail,
			Value: "other@khepri.dev",
		})

		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.user.Id})
		require.NoError(err)

		_, err = other.client.InviteUser(other.ctx, &pb.InviteUserReq{
			OrgId: org.Id[:],
			Identity: &pb.Identity{
				Kind:  pb.IdentityKind_IDENTITY_KIND_MAIL,
				Value: h.identity.Value,
			},
		})
		require.NoError(err)

		_, err = h.client.JoinOrg(ctx, &pb.JoinOrgReq{OrgId: org.Id[:]})
		require.NoError(err)

		member, err := h.Members().GetByUserIdFromOrg(ctx, org.Id, h.user.Id)
		require.NoError(err)
		require.Equal(horus.RoleOrgMember, member.Role)
	}))
}

func TestLeaveOrg(t *testing.T) {
	t.Run("not a member", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		_, err = h.client.LeaveOrg(ctx, &pb.LeaveOrgReq{OrgId: org.Id[:]})
		require.NoError(err)
	}))

	for _, tC := range []struct {
		desc string
		role horus.RoleOrg
	}{
		{
			desc: "as a member",
			role: horus.RoleOrgMember,
		},
		{
			desc: "as a invitee",
			role: horus.RoleOrgInvitee,
		},
	} {
		t.Run(tC.desc, WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
			other, err := h.Users().New(ctx)
			require.NoError(err)

			org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
			require.NoError(err)

			_, err = h.Members().New(ctx, horus.MemberInit{
				OrgId:  org.Id,
				UserId: h.user.Id,
				Role:   tC.role,
			})
			require.NoError(err)

			_, err = h.client.LeaveOrg(ctx, &pb.LeaveOrgReq{OrgId: org.Id[:]})
			require.NoError(err)
		}))
	}

	t.Run("as a member", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		_, err = h.Members().New(ctx, horus.MemberInit{
			OrgId:  org.Id,
			UserId: h.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		_, err = h.client.LeaveOrg(ctx, &pb.LeaveOrgReq{OrgId: org.Id[:]})
		require.NoError(err)
	}))

	t.Run("as an owner", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		_, err = h.Members().New(ctx, horus.MemberInit{
			OrgId:  org.Id,
			UserId: h.user.Id,
			Role:   horus.RoleOrgOwner,
		})
		require.NoError(err)

		_, err = h.client.LeaveOrg(ctx, &pb.LeaveOrgReq{OrgId: org.Id[:]})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.FailedPrecondition, s.Code())
	}))
}

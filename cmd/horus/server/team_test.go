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

func TestNewTeam(t *testing.T) {
	t.Run("as an org owner", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		res, err := h.client.NewTeam(ctx, &pb.NewTeamReq{
			OrgId: rst.Org.Id[:],
			Name:  "A-team",
		})
		require.NoError(err)
		require.Equal("A-team", res.Team.Name)
	}))

	t.Run("as an org member", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		org, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		_, err = h.Members().New(ctx, horus.MemberInit{
			OrgId:  org.Org.Id,
			UserId: h.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		_, err = h.client.NewTeam(ctx, &pb.NewTeamReq{OrgId: org.Org.Id[:]})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.PermissionDenied, s.Code())
	}))

	t.Run("org does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		id := uuid.New()
		_, err := h.client.NewTeam(ctx, &pb.NewTeamReq{
			OrgId: id[:],
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))
}

func TestListTeams(t *testing.T) {
	t.Run("org does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		id := uuid.New()
		_, err := h.client.ListTeams(ctx, &pb.ListTeamsReq{OrgId: id[:]})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("org where the user belongs", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		team1, err := h.Teams().New(ctx, horus.TeamInit{OrgId: rst.Org.Id, OwnerId: rst.Owner.Id})
		require.NoError(err)

		team2, err := h.Teams().New(ctx, horus.TeamInit{OrgId: rst.Org.Id, OwnerId: rst.Owner.Id})
		require.NoError(err)

		res, err := h.client.ListTeams(ctx, &pb.ListTeamsReq{OrgId: rst.Org.Id[:]})
		require.NoError(err)
		require.ElementsMatch([][]byte{team1.Id[:], team2.Id[:]}, fx.MapV(res.Teams, func(v *pb.Team) []byte { return v.Id }))
	}))

	t.Run("org where the user does not belong", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		_, err = h.Teams().New(ctx, horus.TeamInit{OrgId: rst.Org.Id, OwnerId: rst.Owner.Id})
		require.NoError(err)

		_, err = h.client.ListTeams(ctx, &pb.ListTeamsReq{OrgId: rst.Org.Id[:]})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))
}

func TestUpdateTeam(t *testing.T) {
	t.Run("team does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		team_id := uuid.New()
		_, err := h.client.UpdateTeam(ctx, &pb.UpdateTeamReq{
			TeamId: team_id[:],
			Name:   "A-team",
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("as an org member without membership", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
			Name:    "A-team",
		})
		require.NoError(err)

		_, err = h.Members().New(ctx, horus.MemberInit{
			OrgId:  rst.Org.Id,
			UserId: h.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		_, err = h.client.UpdateTeam(ctx, &pb.UpdateTeamReq{
			TeamId: team.Id[:],
			Name:   "Crazy 88s",
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.PermissionDenied, s.Code())
	}))

	t.Run("as not an org member", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
			Name:    "A-team",
		})
		require.NoError(err)

		_, err = h.client.UpdateTeam(ctx, &pb.UpdateTeamReq{
			TeamId: team.Id[:],
			Name:   "Crazy 88s",
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("as a team member", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
			Name:    "A-team",
		})
		require.NoError(err)

		member, err := h.Members().New(ctx, horus.MemberInit{
			OrgId:  rst.Org.Id,
			UserId: h.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		_, err = h.Memberships().New(ctx, horus.MembershipInit{
			TeamId:   team.Id,
			MemberId: member.Id,
			Role:     horus.RoleTeamMember,
		})
		require.NoError(err)

		_, err = h.client.UpdateTeam(ctx, &pb.UpdateTeamReq{
			TeamId: team.Id[:],
			Name:   "Crazy 88s",
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.PermissionDenied, s.Code())
	}))

	t.Run("as a team owner", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
			Name:    "A-team",
		})
		require.NoError(err)

		_, err = h.client.UpdateTeam(ctx, &pb.UpdateTeamReq{
			TeamId: team.Id[:],
			Name:   "Crazy 88s",
		})
		require.NoError(err)
	}))

	t.Run("as an org owner", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
			Name:    "A-team",
		})
		require.NoError(err)

		_, err = h.Members().New(ctx, horus.MemberInit{
			OrgId:  rst.Org.Id,
			UserId: h.user.Id,
			Role:   horus.RoleOrgOwner,
		})
		require.NoError(err)

		_, err = h.client.UpdateTeam(ctx, &pb.UpdateTeamReq{
			TeamId: team.Id[:],
			Name:   "Crazy 88s",
		})
		require.NoError(err)

		team, err = h.Teams().GetById(ctx, team.Id)
		require.NoError(err)
		require.Equal("Crazy 88s", team.Name)
	}))
}

func TestInviteMember(t *testing.T) {
	t.Run("team does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		id := uuid.New()
		_, err = h.client.InviteMember(ctx, &pb.InviteMemberReq{
			TeamId:   id[:],
			MemberId: rst.Owner.Id[:],
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("member does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		id := uuid.New()
		_, err = h.client.InviteMember(ctx, &pb.InviteMemberReq{
			TeamId:   team.Id[:],
			MemberId: id[:],
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("as an org owner without a membership", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		err = h.Memberships().DeleteByUserIdFromTeam(ctx, team.Id, h.user.Id)
		require.NoError(err)

		_, err = h.client.InviteMember(ctx, &pb.InviteMemberReq{
			TeamId:   team.Id[:],
			MemberId: rst.Owner.Id[:],
		})
		require.NoError(err)
	}))

	t.Run("as an org member without a membership", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		member, err := h.Members().New(ctx, horus.MemberInit{
			OrgId:  rst.Org.Id,
			UserId: h.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		_, err = h.client.InviteMember(ctx, &pb.InviteMemberReq{
			TeamId:   team.Id[:],
			MemberId: member.Id[:],
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.PermissionDenied, s.Code())
	}))

	t.Run("as a team owner, invite a member in another org", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst2, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		other_member, err := h.Members().New(ctx, horus.MemberInit{
			OrgId:  rst2.Org.Id,
			UserId: other.Id,
			Role:   horus.RoleOrgOwner,
		})
		require.NoError(err)

		_, err = h.client.InviteMember(ctx, &pb.InviteMemberReq{
			TeamId:   team.Id[:],
			MemberId: other_member.Id[:],
		})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))
}

func TestJoinTeam(t *testing.T) {
	t.Run("team does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		id := uuid.New()
		_, err := h.client.JoinTeam(ctx, &pb.JoinTeamReq{TeamId: id[:]})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("as an org owner", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		other, err := h.Users().New(ctx)
		require.NoError(err)

		other_member, err := h.Members().New(ctx, horus.MemberInit{
			OrgId:  rst.Org.Id,
			UserId: other.Id,
			Role:   horus.RoleOrgOwner,
		})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: other_member.Id,
		})
		require.NoError(err)

		// Note that owner does not need invitation.
		_, err = h.client.JoinTeam(ctx, &pb.JoinTeamReq{TeamId: team.Id[:]})
		require.NoError(err)

		membership, err := h.Memberships().GetByUserIdFromTeam(ctx, team.Id, h.user.Id)
		require.NoError(err)
		require.Equal(horus.RoleTeamMember, membership.Role)
	}))

	t.Run("as an org member without invitation", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		_, err = h.Members().New(ctx, horus.MemberInit{
			OrgId:  rst.Owner.OrgId,
			UserId: h.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		_, err = h.client.JoinTeam(ctx, &pb.JoinTeamReq{TeamId: team.Id[:]})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.PermissionDenied, s.Code())
	}))

	t.Run("as an org member with invitation", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		member, err := h.Members().New(ctx, horus.MemberInit{
			OrgId:  rst.Owner.OrgId,
			UserId: h.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		_, err = h.Memberships().New(ctx, horus.MembershipInit{
			TeamId:   team.Id,
			MemberId: member.Id,
			Role:     horus.RoleTeamInvitee,
		})
		require.NoError(err)

		_, err = h.client.JoinTeam(ctx, &pb.JoinTeamReq{TeamId: team.Id[:]})
		require.NoError(err)
	}))

	t.Run("as not an org member", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		_, err = h.client.JoinTeam(ctx, &pb.JoinTeamReq{TeamId: team.Id[:]})
		s, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, s.Code())
	}))

	t.Run("as a team member", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		_, err = h.client.JoinTeam(ctx, &pb.JoinTeamReq{TeamId: team.Id[:]})
		require.NoError(err)
	}))
}

func TestLeaveTeam(t *testing.T) {
	t.Run("team does not exist", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		id := uuid.New()
		_, err := h.client.LeaveTeam(ctx, &pb.LeaveTeamReq{TeamId: id[:]})
		require.NoError(err)
	}))

	t.Run("as a team owner", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: h.user.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		_, err = h.client.LeaveTeam(ctx, &pb.LeaveTeamReq{TeamId: team.Id[:]})
		require.NoError(err)

		_, err = h.Memberships().GetByUserIdFromTeam(ctx, team.Id, h.user.Id)
		require.ErrorIs(err, horus.ErrNotExist)
	}))

	t.Run("as an org member without a membership", WithHorusGrpc(func(require *require.Assertions, ctx context.Context, h *horusGrpc) {
		other, err := h.Users().New(ctx)
		require.NoError(err)

		rst, err := h.Orgs().New(ctx, horus.OrgInit{OwnerId: other.Id})
		require.NoError(err)

		team, err := h.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		_, err = h.Members().New(ctx, horus.MemberInit{
			OrgId:  rst.Org.Id,
			UserId: h.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		_, err = h.client.LeaveTeam(ctx, &pb.LeaveTeamReq{TeamId: team.Id[:]})
		require.NoError(err)
	}))
}

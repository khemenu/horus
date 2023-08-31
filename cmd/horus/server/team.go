package server

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/pb"
)

func toPbTeam(v *horus.Team) *pb.Team {
	return &pb.Team{
		Id:   v.Id[:],
		Name: v.Name,
	}
}

func fromPbTeam(v *pb.Team) *horus.Team {
	if v == nil {
		return nil
	}

	return &horus.Team{
		Id:   horus.TeamId(v.Id),
		Name: v.Name,
	}
}

func (s *grpcServer) myOrg(ctx context.Context, org_id horus.OrgId) (*horus.Member, error) {
	user := s.mustUser(ctx)

	member, err := s.Members().GetByUserIdFromOrg(ctx, org_id, user.Id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, grpcStatusWithCode(codes.NotFound)
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get member details: %w", err))
	}
	if member.Role != horus.RoleOrgOwner {
		return nil, grpcStatusWithCode(codes.PermissionDenied)
	}

	return member, nil
}

func (s *grpcServer) NewTeam(ctx context.Context, req *pb.NewTeamReq) (*pb.NewTeamRes, error) {
	org_id, err := parseOrgId(req.OrgId)
	if err != nil {
		return nil, err
	}

	owner, err := s.myOrg(ctx, org_id)
	if err != nil {
		return nil, err
	}

	team, err := s.Horus.Teams().New(ctx, horus.TeamInit{
		OrgId:   org_id,
		OwnerId: owner.Id,
		Name:    req.Name,
	})
	if err != nil {
		return nil, grpcInternalErr(ctx, err)
	}

	return &pb.NewTeamRes{
		Team: toPbTeam(team),
	}, nil
}

func (s *grpcServer) ListTeams(ctx context.Context, req *pb.ListTeamsReq) (*pb.ListTeamsRes, error) {
	user := s.mustUser(ctx)

	org_id, err := parseOrgId(req.OrgId)
	if err != nil {
		return nil, err
	}

	_, err = s.Members().GetByUserIdFromOrg(ctx, org_id, user.Id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, grpcStatusWithCode(codes.NotFound)
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get member details: %w", err))
	}

	teams, err := s.Teams().GetAllByOrgId(ctx, org_id)
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("get teams: %w", err))
	}

	return &pb.ListTeamsRes{
		Teams: fx.MapV(teams, toPbTeam),
	}, nil
}

func (s *grpcServer) hasOwnershipOfTeam(ctx context.Context, team_id horus.TeamId) error {
	user := s.mustUser(ctx)

	team, err := s.Teams().GetById(ctx, team_id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return grpcStatusWithCode(codes.NotFound)
		}

		return grpcInternalErr(ctx, fmt.Errorf("get team details: %w", err))
	}

	member, err := s.Members().GetByUserIdFromOrg(ctx, team.OrgId, user.Id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return grpcStatusWithCode(codes.NotFound)
		}

		return grpcInternalErr(ctx, fmt.Errorf("get member details: %w", err))
	}

	if member.Role != horus.RoleOrgOwner {
		membership, err := s.Memberships().GetById(ctx, team.Id, member.Id)
		if err != nil {
			if errors.Is(err, horus.ErrNotExist) {
				return grpcStatusWithCode(codes.PermissionDenied)
			}

			return grpcInternalErr(ctx, fmt.Errorf("get membership details: %w", err))
		}
		if membership.Role != horus.RoleTeamOwner {
			return grpcStatusWithCode(codes.PermissionDenied)
		}
	}

	return nil
}

func (s *grpcServer) UpdateTeam(ctx context.Context, req *pb.UpdateTeamReq) (*pb.UpdateTeamRes, error) {
	team_id, err := parseTeamId(req.TeamId)
	if err != nil {
		return nil, err
	}

	if err := s.hasOwnershipOfTeam(ctx, team_id); err != nil {
		return nil, err
	}

	_, err = s.Teams().UpdateById(ctx, &horus.Team{
		Id:   team_id,
		Name: req.Name,
	})
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("svae: %w", err))
	}

	return &pb.UpdateTeamRes{}, nil
}

func (s *grpcServer) InviteMember(ctx context.Context, req *pb.InviteMemberReq) (*pb.InviteMemberRes, error) {
	team_id, err := parseTeamId(req.TeamId)
	if err != nil {
		return nil, err
	}

	member_id, err := parseMemberId(req.MemberId)
	if err != nil {
		return nil, err
	}

	err = s.hasOwnershipOfTeam(ctx, team_id)
	if err != nil {
		return nil, err
	}

	_, err = s.Memberships().New(ctx, horus.MembershipInit{
		TeamId:   team_id,
		MemberId: member_id,
		Role:     horus.RoleTeamInvitee,
	})
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, grpcStatusWithCode(codes.NotFound)
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("create a membership: %w", err))
	}

	return &pb.InviteMemberRes{}, nil
}

func (s *grpcServer) JoinTeam(ctx context.Context, req *pb.JoinTeamReq) (*pb.JoinTeamRes, error) {
	team_id, err := parseTeamId(req.TeamId)
	if err != nil {
		return nil, err
	}

	user := s.mustUser(ctx)
	membership, err := s.Memberships().GetByUserIdFromTeam(ctx, team_id, user.Id)
	if err == nil {
		if membership.Role != horus.RoleTeamInvitee {
			// Already a member
			return &pb.JoinTeamRes{}, nil
		} else {
			membership.Role = horus.RoleTeamMember
			_, err = s.Memberships().UpdateById(ctx, membership)
			if err != nil {
				return nil, grpcInternalErr(ctx, fmt.Errorf("update membership: %w", err))
			}

			return &pb.JoinTeamRes{}, nil
		}

		// unreachable
	}
	if !errors.Is(err, horus.ErrNotExist) {
		return nil, grpcInternalErr(ctx, fmt.Errorf("get membership details: %w", err))
	}

	member, err := s.Members().GetByUserIdFromTeam(ctx, team_id, user.Id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, grpcStatusWithCode(codes.NotFound)
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get member details: %w", err))
	}
	if member.Role != horus.RoleOrgOwner {
		return nil, grpcStatusWithCode(codes.PermissionDenied)
	}

	_, err = s.Memberships().New(ctx, horus.MembershipInit{
		TeamId:   team_id,
		MemberId: member.Id,
		Role:     horus.RoleTeamMember,
	})
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("create membership: %w", err))
	}

	return &pb.JoinTeamRes{}, nil
}

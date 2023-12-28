package server

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/pb"
)

func toPbMembership(v *horus.Membership) *pb.Membership {
	m := &pb.Membership{
		Team:     toPbTeam(v.Team),
		MemberId: v.MemberId[:],
		Role:     toPbRoleTeam(v.Role),
	}

	if v.Team == nil {
		m.Team = &pb.Team{
			Id: v.TeamId[:],
		}
	}

	return m
}

func (s *grpcServer) GetProfile(ctx context.Context, req *pb.GetProfileReq) (*pb.GetProfileRes, error) {
	var org_id horus.SiloId
	switch k := req.OrgId.(type) {
	case *pb.GetProfileReq_Value:
		var err error
		org_id, err = parseOrgId(k.Value)
		if err != nil {
			return nil, err
		}

	case *pb.GetProfileReq_Alias:
		org, err := s.Orgs().GetByAlias(ctx, k.Alias)
		if err != nil {
			if errors.Is(err, horus.ErrNotExist) {
				return nil, status.Error(codes.NotFound, "not a member")
			}

			return nil, grpcInternalErr(ctx, fmt.Errorf("get org details: %w", err))
		}

		org_id = org.Id
	}

	user := s.mustUser(ctx)
	member, err := s.Members().GetByUserIdFromOrg(ctx, org_id, user.Id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, status.Error(codes.NotFound, "not a member")
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get member details: %w", err))
	}

	memberships, err := s.Memberships().GetAllByMemberId(ctx, member.Id)
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("get membership details: %w", err))
	}

	return &pb.GetProfileRes{
		Profile: &pb.Profile{
			Member:      toPbMember(member),
			Memberships: fx.MapV(memberships, toPbMembership),
		},
	}, nil
}

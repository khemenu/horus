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
	return &pb.Membership{
		TeamId:   v.TeamId[:],
		MemberId: v.MemberId[:],
		Role:     toPbRoleTeam(v.Role),
	}
}

func (s *grpcServer) GetProfile(ctx context.Context, req *pb.GetProfileReq) (*pb.GetProfileRes, error) {
	org_id, err := parseOrgId(req.OrgId)
	if err != nil {
		return nil, err
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

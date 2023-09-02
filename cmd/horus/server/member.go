package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/maps"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/pb"
)

func toPbMember(v *horus.Member) *pb.Member {
	return &pb.Member{
		Id:     v.Id[:],
		OrgId:  v.OrgId[:],
		UserId: v.UserId[:],
		Role:   toPbRoleOrg(v.Role),

		Name:       v.Name,
		Identities: fx.MapV(maps.Values(v.Identities), toPbIdentity),

		CreatedAt: v.CreatedAt.Format(time.RFC3339),
	}
}

func (s *grpcServer) UpdateMember(ctx context.Context, req *pb.UpdateMemberReq) (*pb.UpdateMemberRes, error) {
	org_id, err := parseOrgId(req.OrgId)
	if err != nil {
		return nil, err
	}

	user := s.mustUser(ctx)
	member, err := s.Members().GetByUserIdFromOrg(ctx, org_id, user.Id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, grpcStatusWithCode(codes.NotFound)
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get member details: %w", err))
	}

	member.Name = req.Name
	_, err = s.Members().UpdateById(ctx, member)
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("update a member: %w", err))
	}

	return &pb.UpdateMemberRes{}, nil
}

func (s *grpcServer) AddMemberIdentity(ctx context.Context, req *pb.AddMemberIdentityReq) (*pb.AddMemberIdentityRes, error) {
	org_id, err := parseOrgId(req.OrgId)
	if err != nil {
		return nil, err
	}

	user := s.mustUser(ctx)
	err = s.Members().AddIdentityByUserIdFromOrg(ctx, org_id, user.Id, horus.IdentityValue(req.IdentityValue))
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, grpcStatusWithCode(codes.NotFound)
		}
		if errors.Is(err, horus.ErrInvalidArgument) {
			return nil, grpcStatusWithCode(codes.InvalidArgument)
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("update member: %w", err))
	}

	return &pb.AddMemberIdentityRes{}, nil
}

func (s *grpcServer) RemoveMemberIdentity(ctx context.Context, req *pb.RemoveMemberIdentityReq) (*pb.RemoveMemberIdentityRes, error) {
	org_id, err := parseOrgId(req.OrgId)
	if err != nil {
		return nil, err
	}

	user := s.mustUser(ctx)
	err = s.Members().RemoveIdentityByUserIdFromOrg(ctx, org_id, user.Id, horus.IdentityValue(req.IdentityValue))
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("update member: %w", err))
	}

	return &pb.RemoveMemberIdentityRes{}, nil
}

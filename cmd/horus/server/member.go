package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/maps"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *grpcServer) myMember(ctx context.Context, member_id horus.MemberId) (*horus.Member, error) {
	user := s.mustUser(ctx)

	member, err := s.Members().GetById(ctx, member_id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, grpcStatusWithCode(codes.NotFound)
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get member details: %w", err))
	}

	if user.Id != member.UserId {
		return nil, status.Error(codes.PermissionDenied, "not yours")
	}

	return member, nil
}

func (s *grpcServer) AddMemberIdentity(ctx context.Context, req *pb.AddMemberIdentityReq) (*pb.AddMemberIdentityRes, error) {
	member_id, err := parseMemberId(req.MemberId)
	if err != nil {
		return nil, err
	}

	_, err = s.myMember(ctx, member_id)
	if err != nil {
		return nil, err
	}

	err = s.Members().AddIdentity(ctx, member_id, horus.IdentityValue(req.IdentityValue))
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, status.Errorf(codes.NotFound, "identity does not exist")
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("update member: %w", err))
	}

	return &pb.AddMemberIdentityRes{}, nil
}

func (s *grpcServer) RemoveMemberIdentity(ctx context.Context, req *pb.RemoveMemberIdentityReq) (*pb.RemoveMemberIdentityRes, error) {
	member_id, err := parseMemberId(req.MemberId)
	if err != nil {
		return nil, err
	}

	_, err = s.myMember(ctx, member_id)
	if err != nil {
		return nil, err
	}

	err = s.Members().RemoveIdentity(ctx, member_id, horus.IdentityValue(req.IdentityValue))
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("update member: %w", err))
	}

	return &pb.RemoveMemberIdentityRes{}, nil
}

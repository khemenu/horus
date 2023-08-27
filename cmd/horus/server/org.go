package server

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/pb"
)

func to_org_(v *horus.Org) *pb.Org {
	return &pb.Org{
		Id:   v.Id[:],
		Name: v.Name,
	}
}

func from_org_(v *pb.Org) *horus.Org {
	return &horus.Org{
		Id:   horus.OrgId(v.Id),
		Name: v.Name,
	}
}

func (s *grpcServer) NewOrg(ctx context.Context, req *pb.NewOrgReq) (*pb.NewOrgRes, error) {
	user := s.mustUser(ctx)

	org, err := s.Horus.Orgs().New(ctx, horus.OrgInit{OwnerId: user.Id})
	if err != nil {
		return nil, grpcInternalErr(ctx, err)
	}

	return &pb.NewOrgRes{
		Org: to_org_(org),
	}, nil
}

func (s *grpcServer) ListOrgs(ctx context.Context, req *pb.ListOrgsReq) (*pb.ListOrgsRes, error) {
	user := s.mustUser(ctx)

	orgs, err := s.Horus.Orgs().GetAllByUserId(ctx, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "get org list")
	}

	return &pb.ListOrgsRes{
		Orgs: fx.MapV(orgs, to_org_),
	}, nil
}

func (s *grpcServer) UpdateOrg(ctx context.Context, req *pb.UpdateOrgReq) (*pb.UpdateOrgRes, error) {
	user := s.mustUser(ctx)

	org_id, err := uuid.FromBytes(req.Org.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid organization ID")
	}

	member, err := s.Horus.Members().GetByUserIdFromOrg(ctx, horus.OrgId(org_id), user.Id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, status.Errorf(codes.NotFound, "member details not found")
		}

		return nil, grpcInternalErr(ctx, err)
	}
	if member.Role != horus.RoleOrgOwner {
		return nil, status.Errorf(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	_, err = s.Horus.Orgs().UpdateById(ctx, from_org_(req.Org))
	if err != nil {
		return nil, grpcInternalErr(ctx, err)
	}

	return &pb.UpdateOrgRes{}, nil
}

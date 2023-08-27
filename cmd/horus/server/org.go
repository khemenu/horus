package server

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/pb"
)

func toPbOrg(v *horus.Org) *pb.Org {
	return &pb.Org{
		Id:   v.Id[:],
		Name: v.Name,
	}
}

func fromPbOrg(v *pb.Org) *horus.Org {
	if v == nil {
		return nil
	}

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
		Org: toPbOrg(org),
	}, nil
}

func (s *grpcServer) ListOrgs(ctx context.Context, req *pb.ListOrgsReq) (*pb.ListOrgsRes, error) {
	user := s.mustUser(ctx)

	orgs, err := s.Horus.Orgs().GetAllByUserId(ctx, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "get org list")
	}

	return &pb.ListOrgsRes{
		Orgs: fx.MapV(orgs, toPbOrg),
	}, nil
}

func (s *grpcServer) UpdateOrg(ctx context.Context, req *pb.UpdateOrgReq) (*pb.UpdateOrgRes, error) {
	user := s.mustUser(ctx)

	org_id, err := uuid.FromBytes(req.Org.GetId())
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

	_, err = s.Horus.Orgs().UpdateById(ctx, fromPbOrg(req.Org))
	if err != nil {
		return nil, grpcInternalErr(ctx, err)
	}

	return &pb.UpdateOrgRes{}, nil
}

func (s *grpcServer) InviteUser(ctx context.Context, req *pb.InviteUserReq) (*pb.InviteUserRes, error) {
	user := s.mustUser(ctx)

	org_id, err := uuid.FromBytes(req.OrgId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid organization ID")
	}
	if req.Identity == nil {
		return nil, status.Errorf(codes.InvalidArgument, "identity required")
	}
	if req.Identity.Kind != pb.IdentityKind_IDENTITY_KIND_MAIL {
		// Only mail is supported on current implementation.
		return nil, status.Errorf(codes.InvalidArgument, "only mail is supported")
	}
	if _, err := mail.ParseAddress(req.Identity.Value); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid mail address")
	}

	member, err := s.Members().GetByUserIdFromOrg(ctx, horus.OrgId(org_id), user.Id)
	if err != nil {
		return nil, grpcStatusWithCode(codes.PermissionDenied)
	}
	if member.Role != horus.RoleOrgOwner {
		return nil, grpcStatusWithCode(codes.PermissionDenied)
	}

	identity, err := s.Identities().GetByValue(ctx, req.Identity.Value)
	if err != nil {
		if !errors.Is(err, horus.ErrNotExist) {
			return nil, grpcInternalErr(ctx, fmt.Errorf("get identity details: %w", err))
		}

		identity, err = s.Identities().New(ctx, &horus.IdentityInit{
			Kind:       horus.IdentityMail,
			Value:      req.Identity.Value,
			VerifiedBy: horus.Unverified,
		})
		if err != nil {
			return nil, grpcInternalErr(ctx, fmt.Errorf("create a identity: %w", err))
		}
	}

	_, err = s.Members().New(ctx, horus.MemberInit{
		OrgId:  horus.OrgId(org_id),
		UserId: identity.OwnerId,
		Role:   horus.RoleOrgInvitee,
	})
	if err != nil {
		if errors.Is(err, horus.ErrExist) {
			return nil, status.Errorf(codes.AlreadyExists, "user already a member of the organization")
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("create a member :%w", err))
	}

	return &pb.InviteUserRes{}, nil
}

func (s *grpcServer) JoinOrg(ctx context.Context, req *pb.JoinOrgReq) (*pb.JoinOrgRes, error) {
	user := s.mustUser(ctx)

	org_id, err := uuid.FromBytes(req.OrgId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid organization ID")
	}

	member, err := s.Members().GetByUserIdFromOrg(ctx, horus.OrgId(org_id), user.Id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, grpcStatusWithCode(codes.PermissionDenied)
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get member details: %w", err))
	}

	if member.Role != horus.RoleOrgInvitee {
		return &pb.JoinOrgRes{}, nil
	}

	// Promote from "Invitee" to "Member".
	member.Role = horus.RoleOrgMember
	_, err = s.Members().UpdateById(ctx, member)
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("update a member: %w", err))
	}

	return &pb.JoinOrgRes{}, nil
}

func (s *grpcServer) LeaveOrg(ctx context.Context, req *pb.LeaveOrgReq) (*pb.LeaveOrgRes, error) {
	user := s.mustUser(ctx)

	org_id, err := uuid.FromBytes(req.OrgId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid organization ID")
	}

	member, err := s.Members().GetByUserIdFromOrg(ctx, horus.OrgId(org_id), user.Id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return &pb.LeaveOrgRes{}, nil
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get member details: %w", err))
	}

	if member.Role == horus.RoleOrgOwner {
		return nil, status.Errorf(codes.FailedPrecondition, "owner cannot leave the organization; demote yourself first")
	}

	err = s.Members().DeleteByUserIdFromOrg(ctx, horus.OrgId(org_id), user.Id)
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("delete a member: %w", err))
	}

	return &pb.LeaveOrgRes{}, nil
}

func (s *grpcServer) SetRoleOrg(ctx context.Context, req *pb.SetRoleOrgReq) (*pb.SetRoleOrgRes, error) {
	user := s.mustUser(ctx)

	member_id, err := uuid.FromBytes(req.MemberId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid member ID")
	}
	if fromPbRoleOrg(req.Role) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid role")
	}

	member, err := s.Members().GetById(ctx, horus.MemberId(member_id))
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, grpcStatusWithCode(codes.NotFound)
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get target member details: %w", err))
	}

	me, err := s.Members().GetByUserIdFromOrg(ctx, member.OrgId, user.Id)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			// Target member does not visible to me, so it is not found.
			return nil, grpcStatusWithCode(codes.NotFound)
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get my member details: %w", err))
	}
	if me.Role != horus.RoleOrgOwner {
		return nil, grpcStatusWithCode(codes.PermissionDenied)
	}

	member.Role = fromPbRoleOrg(req.GetRole())
	_, err = s.Members().UpdateById(ctx, member)
	if err != nil {
		if errors.Is(err, horus.ErrFailedPrecondition) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("update a member: %w", err))
	}

	return &pb.SetRoleOrgRes{}, nil
}

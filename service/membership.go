package service

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/proto/khepri/horus"
	"khepri.dev/horus/ent/team"
	"khepri.dev/horus/service/frame"
)

type MembershipService struct {
	horus.UnimplementedMembershipServiceServer
	*base
}

func (s *MembershipService) Create(ctx context.Context, req *horus.CreateMembershipRequest) (*horus.Membership, error) {
	account_id := req.GetMembership().GetAccount().GetId()
	if account_id == nil {
		return nil, newErrMissingRequiredField("membership.account.id")
	}

	team_id := req.GetMembership().GetTeam().GetId()
	if team_id == nil {
		return nil, newErrMissingRequiredField("membership.team.id")
	}

	target_account, err := s.Account().Get(ctx, &horus.GetAccountRequest{
		Id:   account_id,
		View: horus.GetAccountRequest_WITH_EDGE_IDS,
	})
	if err != nil {
		return nil, err
	}

	// Ensure that team exists.
	_, err = s.Team().Get(ctx, &horus.GetTeamRequest{
		Id:   team_id,
		View: horus.GetTeamRequest_WITH_EDGE_IDS,
	})
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if f.ActingAccount == nil {
		// Trying to make myself as a team member.
		if target_account.Role != horus.Account_ROLE_OWNER {
			// Maybe I'm already a member but returns PermissionDenied.
			return nil, ErrPermissionDenied
		}

		return s.store.Membership().Create(ctx, req)
	}
	if f.ActingAccount.Role == account.RoleOWNER {
		if target_account.Role == horus.Account_ROLE_OWNER {
			// Cannot put other owner to the team even if the actor is an owner.
			return nil, ErrPermissionDenied
		}

		return s.store.Membership().Create(ctx, req)
	}

	v, err := f.ActingAccount.QueryMemberships().
		Where(membership.HasTeamWith(team.ID(uuid.UUID(team_id)))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// We know that the team exists, so membership does not exist.
			return nil, ErrPermissionDenied
		}

		return nil, status.Errorf(codes.Internal, "query membership: %s", err.Error())
	}
	if v.Role != membership.RoleOWNER {
		return nil, ErrPermissionDenied
	}

	return s.store.Membership().Create(ctx, req)
}

func (s *MembershipService) Get(ctx context.Context, req *horus.GetMembershipRequest) (*horus.Membership, error) {
	res, err := s.store.Membership().Get(ctx, &horus.GetMembershipRequest{
		Id:   req.Id,
		View: horus.GetMembershipRequest_WITH_EDGE_IDS,
	})
	if err != nil {
		return nil, err
	}
	if _, err := s.Account().Get(ctx, &horus.GetAccountRequest{Id: res.Account.Id}); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *MembershipService) Update(ctx context.Context, req *horus.UpdateMembershipRequest) (*horus.Membership, error) {
	res, err := s.Get(ctx, &horus.GetMembershipRequest{Id: req.Membership.GetId()})
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if f.ActingAccount == nil {
		// Update my self.
		// Only owner of the silo or the team can update the membership.
		return nil, ErrPermissionDenied
	}

	if f.ActingAccount.Role != account.RoleOWNER {
		v, err := f.ActingAccount.QueryMemberships().
			Where(membership.HasTeamWith(team.ID(uuid.UUID(res.Team.Id)))).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, ErrPermissionDenied
			}

			return nil, status.Errorf(codes.Internal, err.Error())
		}
		if v.Role != membership.RoleOWNER {
			return nil, ErrPermissionDenied
		}
	}

	res.Role = req.Membership.Role
	return s.store.Membership().Update(ctx, &horus.UpdateMembershipRequest{
		Membership: res,
	})
}

func (s *MembershipService) Delete(ctx context.Context, req *horus.DeleteMembershipRequest) (*emptypb.Empty, error) {
	res, err := s.Get(ctx, &horus.GetMembershipRequest{Id: req.GetId()})
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if f.ActingAccount == nil {
		return s.store.Membership().Delete(ctx, req)
	}

	if f.ActingAccount.Role != account.RoleOWNER {
		v, err := f.ActingAccount.QueryMemberships().
			Where(membership.HasTeamWith(team.ID(uuid.UUID(res.Team.Id)))).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, ErrPermissionDenied
			}

			return nil, status.Errorf(codes.Internal, err.Error())
		}
		if v.Role != membership.RoleOWNER {
			return nil, ErrPermissionDenied
		}
	}

	return s.store.Membership().Delete(ctx, req)
}

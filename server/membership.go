package server

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/team"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/frame"
)

type MembershipServiceServer struct {
	horus.UnimplementedMembershipServiceServer
	*base
}

func (s *MembershipServiceServer) Create(ctx context.Context, req *horus.CreateMembershipRequest) (*horus.Membership, error) {
	p_acct := req.GetAccount()
	if p_acct == nil {
		return nil, newErrMissingRequiredField(".account.id")
	}

	p_team := req.GetTeam()
	if p_team == nil {
		return nil, newErrMissingRequiredField(".team.id")
	}

	target_account, err := s.bare.Account().Get(ctx, &horus.GetAccountRequest{
		Id: p_acct.Id,
	})
	if err != nil {
		return nil, err
	}

	// Ensure that team exists.
	target_team, err := s.bare.Team().Get(ctx, &horus.GetTeamRequest{
		Id: p_team.Id,
	})
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if f.ActingAccount == nil {
		// Trying to make myself as a team member.
		if target_account.Role != horus.Role_ROLE_OWNER {
			// Maybe I'm already a member but returns PermissionDenied.
			return nil, ErrPermissionDenied
		}

		return s.bare.Membership().Create(ctx, req)
	}
	if f.ActingAccount.Role == role.Owner {
		if target_account.Role == horus.Role_ROLE_OWNER {
			// Cannot put other owner to the team even if the actor is an owner.
			return nil, ErrPermissionDenied
		}

		return s.bare.Membership().Create(ctx, req)
	}

	v, err := f.ActingAccount.QueryMemberships().
		Where(membership.HasTeamWith(team.IDEQ(uuid.UUID(target_team.Id)))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// We know that the team exists, so membership does not exist.
			return nil, ErrPermissionDenied
		}

		return nil, status.Errorf(codes.Internal, "query membership: %s", err.Error())
	}
	if v.Role != role.Owner {
		return nil, ErrPermissionDenied
	}

	return s.bare.Membership().Create(ctx, req)
}

func (s *MembershipServiceServer) Get(ctx context.Context, req *horus.GetMembershipRequest) (*horus.Membership, error) {
	res, err := s.bare.Membership().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	_, err = s.covered.Account().Get(ctx, &horus.GetAccountRequest{
		Id: res.Account.Id,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *MembershipServiceServer) Update(ctx context.Context, req *horus.UpdateMembershipRequest) (*horus.Membership, error) {
	res, err := s.Get(ctx, &horus.GetMembershipRequest{Id: req.GetId()})
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if f.ActingAccount == nil {
		// Update my self.
		// Only owner of the silo or the team can update the membership.
		return nil, ErrPermissionDenied
	}

	if f.ActingAccount.Role != role.Owner {
		v, err := f.ActingAccount.QueryMemberships().
			Where(membership.HasTeamWith(team.ID(uuid.UUID(res.Team.Id)))).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, ErrPermissionDenied
			}

			return nil, status.Errorf(codes.Internal, err.Error())
		}
		if v.Role != role.Owner {
			return nil, ErrPermissionDenied
		}
	}

	return s.bare.Membership().Update(ctx, req)
}

func (s *MembershipServiceServer) Delete(ctx context.Context, req *horus.DeleteMembershipRequest) (*emptypb.Empty, error) {
	res, err := s.Get(ctx, &horus.GetMembershipRequest{Id: req.GetId()})
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if f.ActingAccount == nil {
		return s.bare.Membership().Delete(ctx, req)
	}

	if f.ActingAccount.Role != role.Owner {
		v, err := f.ActingAccount.QueryMemberships().
			Where(membership.HasTeamWith(team.ID(uuid.UUID(res.Team.Id)))).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, ErrPermissionDenied
			}

			return nil, status.Errorf(codes.Internal, err.Error())
		}
		if v.Role != role.Owner {
			return nil, ErrPermissionDenied
		}
	}

	return s.bare.Membership().Delete(ctx, req)
}

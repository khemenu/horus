package server

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/predicate"
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/ent/team"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/bare"
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

func (s *MembershipServiceServer) List(ctx context.Context, req *horus.ListMembershipRequest) (*horus.ListMembershipResponse, error) {
	l := int(req.GetLimit())
	l = fx.Clamp(l, 5, 100)
	q := s.db.Membership.Query().
		Order(membership.ByDateCreated(sql.OrderDesc())).
		Limit(l)

	ps := make([]predicate.Membership, 0, 2)
	if t := req.GetToken(); t != nil {
		ps = append(ps, membership.DateCreatedLT(t.AsTime()))
	}

	var (
		vs  []*ent.Membership
		err error
	)
	r := &horus.GetSiloRequest{}
	switch k := req.GetKey().(type) {
	case *horus.ListMembershipRequest_Mine:
		f := frame.Must(ctx)
		ps = append(ps, membership.HasAccountWith(
			account.HasOwnerWith(user.IDEQ(f.Actor.ID)),
		))
		vs, err = q.Where(ps...).
			WithAccount().
			WithTeam().
			All(ctx)
		if err != nil {
			return nil, bare.ToStatus(err)
		}
		goto R

	case *horus.ListMembershipRequest_SiloId:
		r.Key = &horus.GetSiloRequest_Id{Id: k.SiloId}
	case *horus.ListMembershipRequest_SiloAlias:
		r.Key = &horus.GetSiloRequest_Alias{Alias: k.SiloAlias}

	case *horus.ListMembershipRequest_TeamId:
		return nil, status.Error(codes.Unimplemented, "not implemented for given key")
	case *horus.ListMembershipRequest_TeamAlias:
		return nil, status.Error(codes.Unimplemented, "not implemented for given key")

	default:
		return nil, status.Error(codes.InvalidArgument, "unknown key")
	}

	if v, err := s.covered.Silo().Get(ctx, r); err != nil {
		return nil, err
	} else {
		ps = append(ps, membership.HasAccountWith(
			account.HasSiloWith(silo.IDEQ(uuid.UUID(v.Id))),
		))
	}

	vs, err = q.Where(ps...).
		WithAccount(func(q *ent.AccountQuery) {
			q.WithOwner()
			q.WithSilo()
		}).
		WithTeam().
		All(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

R:
	return &horus.ListMembershipResponse{
		Items: fx.MapV(vs, func(v *ent.Membership) *horus.Membership {
			return bare.ToProtoMembership(v)
		}),
	}, nil
}

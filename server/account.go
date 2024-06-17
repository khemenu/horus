package server

import (
	"bytes"
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/predicate"
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

type AccountServiceServer struct {
	horus.UnimplementedAccountServiceServer
	*base
}

func (s *AccountServiceServer) Create(ctx context.Context, req *horus.CreateAccountRequest) (*horus.Account, error) {
	f := frame.Must(ctx)

	owner := req.GetOwner()
	if owner == nil {
		return nil, status.Errorf(codes.PermissionDenied, "account cannot be created by the account holder themselves")
	}

	// Actor must be owner of the silo.
	// Actor must be direct parent of the account owner.
	silo_uuid, err := bare.GetSiloId(ctx, s.db, req.GetSilo())
	if err != nil {
		return nil, err
	}

	if actor_acct, err := s.db.Account.Query().
		Where(
			account.OwnerIDEQ(f.Actor.ID),
			account.SiloIDEQ(silo_uuid),
		).
		Only(ctx); err != nil || actor_acct.Role != role.Owner {
		return nil, status.Error(codes.PermissionDenied, "not the silo owner")
	}

	var owner_uuid uuid.UUID
	q := s.db.User.Query().WithParent()
	if p, err := bare.GetUserSpecifier(req.GetOwner()); err != nil {
		return nil, err
	} else {
		q.Where(p)
	}
	if owner, err := q.Only(ctx); err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "owner not found")
		}

		return nil, fmt.Errorf("get silo: %w", err)
	} else if p := owner.Edges.Parent; p == nil || p.ID != f.Actor.ID {
		return nil, status.Error(codes.PermissionDenied, "actor is not a direct parent of the account owner")
	} else {
		owner_uuid = owner.ID
	}

	req.Role = horus.Role_ROLE_MEMBER
	req.Owner = horus.UserById(owner_uuid)
	req.Silo = horus.SiloById(silo_uuid)
	return s.bare.Account().Create(ctx, req)
}

func (s *AccountServiceServer) Get(ctx context.Context, req *horus.GetAccountRequest) (*horus.Account, error) {
	res, err := s.bare.Account().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if bytes.Equal(f.Actor.ID[:], res.Owner.Id) {
		return res, nil
	}

	v, err := f.Actor.QueryAccounts().
		Where(account.SiloID(uuid.UUID(res.Silo.Id))).
		Only(ctx)
	if err == nil {
		f.ActingAccount = v
		return res, nil
	}
	if ent.IsNotFound(err) {
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	}

	return nil, status.Errorf(codes.Internal, err.Error())
}

func (s *AccountServiceServer) Update(ctx context.Context, req *horus.UpdateAccountRequest) (*horus.Account, error) {
	v, err := s.Get(ctx, req.GetKey())
	if err != nil {
		return nil, err
	}
	if v.Role != horus.Role_ROLE_OWNER {
		req.Role = nil
	}

	f := frame.Must(ctx)
	if !bytes.Equal(f.Actor.ID[:], v.Owner.Id) {
		return nil, ErrPermissionDenied
	}

	return s.bare.Account().Update(ctx, req)
}

func (s *AccountServiceServer) Delete(ctx context.Context, req *horus.GetAccountRequest) (*emptypb.Empty, error) {
	f := frame.Must(ctx)

	v, err := s.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	if v.Role == horus.Role_ROLE_OWNER {
		return nil, status.Error(codes.PermissionDenied, "owner account cannot be deleted manually")
	}
	switch {
	case bytes.Equal(f.Actor.ID[:], v.Owner.Id):
		// Delete myself.
		fallthrough
	case f.ActingAccount.Role == role.Owner:
		// Delete account in my own silo.
		return s.bare.Account().Delete(ctx, req)
	}

	return nil, ErrPermissionDenied
}

func (s *AccountServiceServer) List(ctx context.Context, req *horus.ListAccountRequest) (*horus.ListAccountResponse, error) {
	q := s.db.Account.Query().
		Order(account.ByDateCreated(sql.OrderDesc()))
	if l := req.GetLimit(); l > 0 {
		q.Limit(int(l))
	}

	ps := make([]predicate.Account, 0, 2)
	if t := req.GetToken(); t != nil {
		ps = append(ps, account.DateCreatedLT(t.AsTime()))
	}

	var (
		vs  []*ent.Account
		err error
	)
	r := &horus.GetSiloRequest{}
	switch k := req.GetKey().(type) {
	case *horus.ListAccountRequest_Mine:
		f := frame.Must(ctx)
		ps = append(ps, account.HasOwnerWith(user.IDEQ(f.Actor.ID)))
		vs, err = q.Where(ps...).
			WithSilo().
			All(ctx)
		if err != nil {
			return nil, bare.ToStatus(err)
		}
		goto R

	case *horus.ListAccountRequest_SiloId:
		r.Key = &horus.GetSiloRequest_Id{Id: k.SiloId}
	case *horus.ListAccountRequest_SiloAlias:
		r.Key = &horus.GetSiloRequest_Alias{Alias: k.SiloAlias}
	default:
		return nil, status.Error(codes.InvalidArgument, "unknown key")
	}

	if v, err := s.covered.Silo().Get(ctx, r); err != nil {
		return nil, err
	} else {
		ps = append(ps, account.HasSiloWith(silo.IDEQ(uuid.UUID(v.Id))))
	}

	vs, err = q.Where(ps...).
		WithMemberships(func(q *ent.MembershipQuery) {
			q.WithTeam()
		}).
		WithOwner().
		All(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

R:
	return &horus.ListAccountResponse{
		Items: fx.MapV(vs, func(v *ent.Account) *horus.Account {
			return bare.ToProtoAccount(v)
		}),
	}, nil
}

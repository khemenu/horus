package server

import (
	"bytes"
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
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

	// Actor must be a silo owner or a silo admin.
	// Actor must be a parent of the account owner.
	// Actor except for the owner must have higher role than creating account.

	if p, err := bare.GetSiloSpecifier(req.GetSilo()); err != nil {
		return nil, err
	} else if actor_account, err := s.db.Account.Query().
		Where(
			account.OwnerIDEQ(f.Actor.ID),
			account.HasSiloWith(p),
		).
		WithSilo(func(q *ent.SiloQuery) { q.Select(silo.FieldID) }).
		Only(ctx); err != nil {
		return nil, bare.ToStatus(err)
	} else if actor_account.Role != role.Owner && actor_account.Role != role.Admin {
		return nil, status.Error(codes.PermissionDenied, "account can be created only by a silo owner or a silo admin")
	} else if actor_account.Role != role.Owner && !req.GetRole().LowerThan(actor_account.Role) {
		return nil, status.Error(codes.PermissionDenied, "account cannot be created if the role is not lower than the role of actor's account")
	} else {
		req.Silo = horus.SiloById(actor_account.SiloID)
	}

	if p, err := bare.GetUserSpecifier(req.GetOwner()); err != nil {
		return nil, err
	} else if owner, err := s.db.User.Query().Where(p).WithParent().Only(ctx); err != nil {
		return nil, bare.ToStatus(err)
	} else if p := owner.Edges.Parent; p == nil || p.ID != f.Actor.ID {
		return nil, status.Error(codes.FailedPrecondition, "account cannot be created for a user who is not your child")
	} else {
		req.Owner = horus.UserById(owner.ID)
	}

	if req.Role == horus.Role_ROLE_UNSPECIFIED {
		req.Role = horus.Role_ROLE_MEMBER
	}
	return s.bare.Account().Create(ctx, req)
}

func (s *AccountServiceServer) Get(ctx context.Context, req *horus.GetAccountRequest) (*horus.Account, error) {
	f := frame.Must(ctx)

	if k := req.GetByAliasInSilo(); k.GetAlias() == horus.Me {
		p, err := bare.GetSiloSpecifier(k.Silo)
		if err != nil {
			return nil, err
		}

		v, err := bare.QueryAccountWithEdgeIds(f.Actor.QueryAccounts()).
			Where(account.HasSiloWith(p)).
			Only(ctx)
		if err != nil {
			return nil, bare.ToStatus(err)
		}

		f.ActingAccount = v
		return bare.ToProtoAccount(v), nil
	}

	v, err := s.bare.Account().Get(ctx, req)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(v.Owner.Id, f.Actor.ID[:]) {
		return v, nil
	}

	// Test if the actor has an account in the same silo where the account to get.
	w, err := bare.QueryAccountWithEdgeIds(f.Actor.QueryAccounts()).
		Where(account.SiloID(uuid.UUID(v.Silo.Id))).
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	f.ActingAccount = w
	return v, nil
}

func (s *AccountServiceServer) Update(ctx context.Context, req *horus.UpdateAccountRequest) (*horus.Account, error) {
	f := frame.Must(ctx)

	v, err := s.Get(ctx, req.GetKey())
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(v.Owner.Id, f.Actor.ID[:]) {
		// Updating the another user's account.
		if my_role := f.MustGetActingAccount().Role; my_role != role.Owner {
			if r := req.GetRole(); r != horus.Role_ROLE_UNSPECIFIED {
				return nil, status.Error(codes.PermissionDenied, "account role can be updated only by the silo owner")
			}
			if !v.Role.LowerThan(my_role) {
				return nil, status.Error(codes.PermissionDenied, "account can be updated only by the account owner or a user with a higher role")
			}
		}
	} else if r := req.GetRole(); r != horus.Role_ROLE_UNSPECIFIED {
		// Updating my account role.
		if r > v.Role {
			return nil, status.Error(codes.PermissionDenied, "cannot promote yourself")
		}
		if r < v.Role && v.Role == horus.Role_ROLE_OWNER {
			// The owner of the silo trying to demote itself.
			n, err := s.db.Account.Query().
				Where(
					account.HasSiloWith(silo.IDEQ(uuid.UUID(v.Silo.Id))),
					account.RoleEQ(role.Owner),
				).
				Count(ctx)
			if err != nil {
				return nil, bare.ToStatus(err)
			}
			if n == 1 {
				return nil, status.Error(codes.FailedPrecondition, "sole owner cannot be demoted")
			}
		}
	}

	return s.bare.Account().Update(ctx, req)
}

func (s *AccountServiceServer) Delete(ctx context.Context, req *horus.GetAccountRequest) (*emptypb.Empty, error) {
	f := frame.Must(ctx)

	v, err := s.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(v.Owner.Id, f.Actor.ID[:]) {
		// Deleting another user's account.
		my_role := f.MustGetActingAccount().Role
		if my_role != role.Owner && !v.Role.LowerThan(my_role) {
			return nil, status.Error(codes.PermissionDenied, "account can be deleted only by the account owner or a user with a higher role")
		}
	} else if v.Role == horus.Role_ROLE_OWNER {
		// The owner of the silo trying to delete itself.
		n, err := s.db.Account.Query().
			Where(
				account.HasSiloWith(silo.IDEQ(uuid.UUID(v.Silo.Id))),
				account.RoleEQ(role.Owner),
			).
			Count(ctx)
		if err != nil {
			return nil, bare.ToStatus(err)
		}
		if n == 1 {
			return nil, status.Error(codes.FailedPrecondition, "sole owner cannot be deleted")
		}
	}

	return s.bare.Account().Delete(ctx, req)
}

func (s *AccountServiceServer) List(ctx context.Context, req *horus.ListAccountRequest) (*horus.ListAccountResponse, error) {
	f := frame.Must(ctx)

	q := s.db.Account.Query().
		Order(account.ByDateCreated(sql.OrderDesc()))
	if l := req.GetLimit(); l > 0 {
		q.Limit(int(l))
	}
	if t := req.GetToken(); t != nil {
		q.Where(account.DateCreatedLT(t.AsTime()))
	}

	var (
		vs  []*ent.Account
		err error
	)
	switch k := req.GetKey().(type) {
	case *horus.ListAccountRequest_Mine:
		vs, err = q.Where(account.HasOwnerWith(user.IDEQ(f.Actor.ID))).WithSilo().All(ctx)
		goto R

	case *horus.ListAccountRequest_Silo:
		v, err := s.covered.Silo().Get(ctx, k.Silo)
		if err != nil {
			return nil, err
		}
		q.Where(account.HasSiloWith(silo.IDEQ(uuid.UUID(v.Id))))

	case nil:
		return nil, status.Error(codes.InvalidArgument, "key not provided")
	default:
		return nil, status.Error(codes.Unimplemented, "unknown key")
	}

	vs, err = q.
		WithMemberships(func(q *ent.MembershipQuery) {
			q.WithTeam()
		}).
		WithOwner().
		All(ctx)

R:
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	return &horus.ListAccountResponse{
		Items: fx.MapV(vs, func(v *ent.Account) *horus.Account {
			return bare.ToProtoAccount(v)
		}),
	}, nil
}

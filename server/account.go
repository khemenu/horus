package server

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lesomnus/entpb/cmd/protoc-gen-entpb/runtime"
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

	target_silo := req.GetSilo()
	silo_uuid, err := uuid.FromBytes(target_silo.GetId())
	if err != nil && target_silo.GetId() != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid silo UUID")
	}
	if silo_uuid == uuid.Nil {
		silo_alias := target_silo.GetAlias()
		if silo_alias == "" {
			return nil, status.Errorf(codes.InvalidArgument, "silo ID not provided")
		}

		silo_uuid, err = s.db.Silo.Query().
			Where(silo.AliasEQ(silo_alias)).
			OnlyID(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, status.Error(codes.NotFound, "silo not found")
			}

			return nil, fmt.Errorf("get silo: %w", err)
		}
	}

	if actor_acct, err := s.db.Account.Query().
		Where(
			account.OwnerIDEQ(f.Actor.ID),
			account.SiloIDEQ(silo_uuid),
		).
		Only(ctx); err != nil || actor_acct.Role != role.Owner {
		return nil, status.Error(codes.PermissionDenied, "not the silo owner")
	}

	owner_uuid, err := uuid.FromBytes(owner.Id)
	if err != nil && owner.Id != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner UUID")
	}

	var owner_q predicate.User
	if owner_uuid != uuid.Nil {
		owner_q = user.IDEQ(owner_uuid)
	} else {
		if owner.Alias == "" {
			return nil, status.Errorf(codes.InvalidArgument, "owner ID not provided")
		}

		owner_q = user.AliasEQ(owner.Alias)
	}

	if owner, err := s.db.User.Query().
		Where(owner_q).
		WithParent().
		Only(ctx); err != nil {
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
	req.Owner = &horus.User{Id: owner_uuid[:]}
	req.Silo = &horus.Silo{Id: silo_uuid[:]}
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
	v, err := s.Get(ctx, &horus.GetAccountRequest{
		Id: req.GetId(),
	})
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

func (s *AccountServiceServer) Delete(ctx context.Context, req *horus.DeleteAccountRequest) (*emptypb.Empty, error) {
	f := frame.Must(ctx)
	v, err := s.Get(ctx, &horus.GetAccountRequest{
		Id: req.GetId(),
	})
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
	f := frame.Must(ctx)
	vs, err := s.db.Account.Query().
		Where(account.OwnerIDEQ(f.Actor.ID)).
		WithSilo().
		All(ctx)
	if err != nil {
		return nil, runtime.EntErrorToStatus(err)
	}

	return &horus.ListAccountResponse{
		Items: fx.MapV(vs, func(v *ent.Account) *horus.Account {
			return bare.ToProtoAccount(v)
		}),
	}, nil
}

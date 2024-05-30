package server

import (
	"bytes"
	"context"
	"fmt"

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
	"khepri.dev/horus/server/frame"
)

type AccountServiceServer struct {
	horus.UnimplementedAccountServiceServer
	*base
}

func (s *AccountServiceServer) Create(ctx context.Context, req *horus.CreateAccountRequest) (*horus.Account, error) {
	f := frame.Must(ctx)

	owner := req.GetAccount().GetOwner()
	if owner == nil {
		return nil, status.Errorf(codes.PermissionDenied, "account cannot be created by the account holder themselves")
	}

	// Actor must be owner of the silo.
	// Actor must be direct parent of the account owner.

	target_silo := req.GetAccount().GetSilo()
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
		Only(ctx); err != nil || actor_acct.Role != account.RoleOWNER {
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

	v := req.GetAccount()
	return s.bare.Account().Create(ctx, &horus.CreateAccountRequest{Account: &horus.Account{
		Alias:       v.GetAlias(),
		Name:        v.GetName(),
		Description: v.GetDescription(),
		Role:        horus.Account_ROLE_MEMBER,

		Owner: &horus.User{Id: owner_uuid[:]},
		Silo:  &horus.Silo{Id: silo_uuid[:]},
	}})
}

func (s *AccountServiceServer) Get(ctx context.Context, req *horus.GetAccountRequest) (*horus.Account, error) {
	res, err := s.bare.Account().Get(ctx, &horus.GetAccountRequest{
		Id:   req.Id,
		View: horus.GetAccountRequest_WITH_EDGE_IDS,
	})
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

func (s *AccountServiceServer) List(ctx context.Context, req *horus.ListAccountRequest) (*horus.ListAccountResponse, error) {
	// f := frame.Must(ctx)
	// vs, err := f.Actor.QueryAccounts().
	// 	WithSilo().
	// 	All(ctx)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, err.Error())
	// }

	// res, err := horus.ToProtoAccountList(vs)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, err.Error())
	// }

	// for i, v := range vs {
	// 	res[i].Silo, err = horus.ToProtoSilo(v.Edges.Silo)
	// 	if err != nil {
	// 		return nil, status.Errorf(codes.Internal, err.Error())
	// 	}
	// }

	// return &horus.ListAccountResponse{AccountList: res}, nil

	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func (s *AccountServiceServer) Update(ctx context.Context, req *horus.UpdateAccountRequest) (*horus.Account, error) {
	v, err := s.Get(ctx, &horus.GetAccountRequest{Id: req.Account.Id})
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if !bytes.Equal(f.Actor.ID[:], v.Owner.Id) {
		return nil, ErrPermissionDenied
	}

	v.Alias = req.Account.Alias
	v.Name = req.Account.Name
	v.Description = req.Account.Description
	return s.bare.Account().Update(ctx, &horus.UpdateAccountRequest{
		Account: v,
	})
}

func (s *AccountServiceServer) Delete(ctx context.Context, req *horus.DeleteAccountRequest) (*emptypb.Empty, error) {
	f := frame.Must(ctx)
	v, err := s.Get(ctx, &horus.GetAccountRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	if v.Role == horus.Account_ROLE_OWNER {
		return nil, status.Error(codes.PermissionDenied, "owner account cannot be deleted manually")
	}
	switch {
	case bytes.Equal(f.Actor.ID[:], v.Owner.Id):
		fallthrough
	case f.ActingAccount.Role == account.RoleOWNER:
		return s.bare.Account().Delete(ctx, req)
	}

	return nil, ErrPermissionDenied
}

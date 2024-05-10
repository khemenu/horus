package service

import (
	"bytes"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/service/frame"
)

type AccountService struct {
	horus.UnimplementedAccountServiceServer
	*base
}

func (s *AccountService) Create(ctx context.Context, req *horus.CreateAccountRequest) (*horus.Account, error) {
	return nil, status.Errorf(codes.PermissionDenied, "account cannot be created manually")
}

func (s *AccountService) Get(ctx context.Context, req *horus.GetAccountRequest) (*horus.Account, error) {
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

func (s *AccountService) List(ctx context.Context, req *horus.ListAccountRequest) (*horus.ListAccountResponse, error) {
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

func (s *AccountService) Update(ctx context.Context, req *horus.UpdateAccountRequest) (*horus.Account, error) {
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

func (s *AccountService) Delete(ctx context.Context, req *horus.DeleteAccountRequest) (*emptypb.Empty, error) {
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

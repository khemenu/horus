package server

import (
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
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/internal/entutils"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

type SiloServiceServer struct {
	horus.UnimplementedSiloServiceServer
	*base
}

func (s *SiloServiceServer) Create(ctx context.Context, req *horus.CreateSiloRequest) (*horus.Silo, error) {
	f := frame.Must(ctx)
	if req == nil {
		req = &horus.CreateSiloRequest{}
	}
	return entutils.WithTxV(ctx, s.db, func(tx *ent.Tx) (*horus.Silo, error) {
		c := tx.Client()
		v, err := bare.NewSiloServiceServer(c).Create(ctx, req)
		if err != nil {
			return nil, err
		}

		_, err = bare.NewAccountServiceServer(c).Create(ctx, &horus.CreateAccountRequest{
			Alias:       fx.Addr("founder"),
			Name:        fx.Addr("Founder"),
			Description: fx.Addr(fmt.Sprintf("Founder of %s", v.Name)),

			Role: horus.Role_ROLE_OWNER,

			Owner: &horus.User{Id: f.Actor.ID[:]},
			Silo:  &horus.Silo{Id: v.Id},
		})
		if err != nil {
			return nil, err
		}

		return v, nil
	})
}

func (s *SiloServiceServer) Get(ctx context.Context, req *horus.GetSiloRequest) (*horus.Silo, error) {
	f := frame.Must(ctx)
	res, err := s.bare.Silo().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	v, err := f.Actor.QueryAccounts().
		Where(account.HasSiloWith(silo.ID(uuid.UUID(res.GetId())))).
		Only(ctx)
	if err != nil {
		return nil, runtime.EntErrorToStatus(err)
	}

	f.ActingAccount = v
	return res, nil
}

func (s *SiloServiceServer) Update(ctx context.Context, req *horus.UpdateSiloRequest) (*horus.Silo, error) {
	_, err := s.Get(ctx, &horus.GetSiloRequest{Key: &horus.GetSiloRequest_Id{
		Id: req.GetId(),
	}})
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if f.ActingAccount.Role != role.Owner {
		return nil, ErrPermissionDenied
	}

	return s.bare.Silo().Update(ctx, req)
}

func (s *SiloServiceServer) Delete(ctx context.Context, req *horus.DeleteSiloRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.PermissionDenied, "silo cannot be deleted manually")
}

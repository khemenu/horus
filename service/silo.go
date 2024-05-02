package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus/alias"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/proto/khepri/horus"
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/internal/entutils"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/service/frame"
)

type SiloService struct {
	horus.UnimplementedSiloServiceServer
	*base
}

func (s *SiloService) Create(ctx context.Context, req *horus.CreateSiloRequest) (*horus.Silo, error) {
	f := frame.Must(ctx)
	return entutils.WithTxV(ctx, s.client, func(tx *ent.Tx) (*horus.Silo, error) {
		c := tx.Client()
		res, err := horus.NewSiloService(c).Create(ctx, &horus.CreateSiloRequest{
			Silo: &horus.Silo{
				Alias: fx.CoalesceOr(req.Silo.GetAlias(), alias.New()),
				Name:  req.Silo.GetName(),
			},
		})
		if err != nil {
			return nil, err
		}

		_, err = horus.NewAccountService(c).Create(ctx, &horus.CreateAccountRequest{
			Account: &horus.Account{
				Alias:       "founder",
				Name:        "Founder",
				Description: fmt.Sprintf("Founder of %s", res.Name),

				Role: horus.Account_ROLE_OWNER,

				Owner: &horus.User{Id: f.Actor.ID[:]},
				Silo:  &horus.Silo{Id: res.Id},
			},
		})
		if err != nil {
			return nil, err
		}

		return res, nil
	})
}

func (s *SiloService) Get(ctx context.Context, req *horus.GetSiloRequest) (*horus.Silo, error) {
	f := frame.Must(ctx)
	res, err := s.store.Silo().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	v, err := f.Actor.QueryAccounts().
		Where(account.HasSiloWith(silo.ID(uuid.UUID(req.Id)))).
		Only(ctx)
	if err == nil {
		f.ActingAccount = v
		return res, nil
	}
	if ent.IsNotFound(err) {
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	}

	return nil, status.Errorf(codes.Internal, "internal error: %s", err)
}

func (s *SiloService) Update(ctx context.Context, req *horus.UpdateSiloRequest) (*horus.Silo, error) {
	v, err := s.Get(ctx, &horus.GetSiloRequest{Id: req.Silo.Id})
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if f.ActingAccount.Role != account.RoleOWNER {
		return nil, ErrPermissionDenied
	}

	v.Alias = req.Silo.Alias
	v.Name = req.Silo.Name
	v.Description = req.Silo.Description
	return s.store.Silo().Update(ctx, &horus.UpdateSiloRequest{
		Silo: v,
	})
}

func (s *SiloService) Delete(ctx context.Context, req *horus.DeleteSiloRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.PermissionDenied, "silo cannot be deleted manually")
}

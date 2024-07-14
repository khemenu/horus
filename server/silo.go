package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
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

			Owner: horus.UserById(f.Actor.ID),
			Silo:  horus.SiloByIdV(v.Id),
		})
		if err != nil {
			return nil, err
		}

		return v, nil
	})
}

func (s *SiloServiceServer) Get(ctx context.Context, req *horus.GetSiloRequest) (*horus.Silo, error) {
	f := frame.Must(ctx)

	p, err := bare.GetSiloSpecifier(req)
	if err != nil {
		return nil, err
	}

	v, err := f.Actor.QueryAccounts().
		Where(account.HasSiloWith(p)).
		WithSilo().
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	f.ActingAccount = v
	return bare.ToProtoSilo(v.Edges.Silo), nil
}

func (s *SiloServiceServer) actorCanModify(ctx context.Context, req *horus.GetSiloRequest) error {
	f := frame.Must(ctx)

	p, err := bare.GetSiloSpecifier(req)
	if err != nil {
		return err
	}

	v, err := f.Actor.QueryAccounts().
		Where(account.HasSiloWith(p)).
		Only(ctx)
	if err != nil {
		return bare.ToStatus(err)
	}
	if v.Role != role.Owner {
		return status.Error(codes.PermissionDenied, "silo can only be modified by its owner")
	}

	return nil
}

func (s *SiloServiceServer) Update(ctx context.Context, req *horus.UpdateSiloRequest) (*horus.Silo, error) {
	if err := s.actorCanModify(ctx, req.GetKey()); err != nil {
		return nil, err
	}

	return s.bare.Silo().Update(ctx, req)
}

func (s *SiloServiceServer) Delete(ctx context.Context, req *horus.GetSiloRequest) (*emptypb.Empty, error) {
	if err := s.actorCanModify(ctx, req); err != nil {
		return nil, err
	}

	return s.bare.Silo().Delete(ctx, req)
}

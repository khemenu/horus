package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/ent/team"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

type ConfServiceServer struct {
	horus.UnimplementedConfServiceServer
	*base
}

func (s *ConfServiceServer) actorCan(ctx context.Context, act string) error {
	f := frame.Must(ctx)
	a, err := f.Actor.QueryAccounts().
		Where(account.HasSiloWith(silo.AliasEQ(horus.ConfSiloName))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return status.Errorf(codes.PermissionDenied, codes.PermissionDenied.String())
		}

		return bare.ToStatus(err)
	}
	if a.Role == role.Owner || a.Role == role.Admin {
		return nil
	}

	_, err = a.QueryMemberships().
		Where(membership.HasTeamWith(team.AliasEQ(act))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return status.Errorf(codes.PermissionDenied, codes.PermissionDenied.String())
		}

		return bare.ToStatus(err)
	}

	return nil
}

func (s *ConfServiceServer) Create(ctx context.Context, req *horus.CreateConfRequest) (*horus.Conf, error) {
	if err := s.actorCan(ctx, "create"); err != nil {
		return nil, err
	}

	return s.bare.Conf().Create(ctx, req)
}

func (s *ConfServiceServer) Delete(ctx context.Context, req *horus.GetConfRequest) (*emptypb.Empty, error) {
	if err := s.actorCan(ctx, "delete"); err != nil {
		return nil, err
	}

	return s.bare.Conf().Delete(ctx, req)
}

func (s *ConfServiceServer) Get(ctx context.Context, req *horus.GetConfRequest) (*horus.Conf, error) {
	if err := s.actorCan(ctx, "get"); err != nil {
		return nil, err
	}

	return s.bare.Conf().Get(ctx, req)
}

func (s *ConfServiceServer) Update(ctx context.Context, req *horus.UpdateConfRequest) (*horus.Conf, error) {
	if err := s.actorCan(ctx, "update"); err != nil {
		return nil, err
	}

	return s.bare.Conf().Update(ctx, req)
}

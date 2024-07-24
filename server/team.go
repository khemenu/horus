package server

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/ent/team"
	"khepri.dev/horus/internal/entutils"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

type TeamServiceServer struct {
	horus.UnimplementedTeamServiceServer
	*base
}

func (s *TeamServiceServer) Create(ctx context.Context, req *horus.CreateTeamRequest) (*horus.Team, error) {
	f := frame.Must(ctx)

	p, err := bare.GetSiloSpecifier(req.GetSilo())
	if err != nil {
		return nil, err
	}

	a, err := f.Actor.QueryAccounts().
		Where(account.HasSiloWith(p)).
		WithSilo().
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}
	switch a.Role {
	case role.Owner:
		fallthrough
	case role.Admin:
		break

	default:
		return nil, status.Error(codes.PermissionDenied, "team can be created only by the silo owners or silo admins")
	}

	req.Silo = horus.SiloById(a.Edges.Silo.ID)
	return entutils.WithTxV(ctx, s.db, func(tx *ent.Tx) (*horus.Team, error) {
		c := tx.Client()
		s, err := bare.NewTeamServiceServer(c).Create(ctx, req)
		if err != nil {
			return nil, err
		}

		_, err = bare.NewMembershipServiceServer(c).Create(ctx, &horus.CreateMembershipRequest{
			Role:    fx.Addr(horus.Role_ROLE_OWNER),
			Account: horus.AccountById(a.ID),
			Team:    horus.TeamByIdV(s.Id),
		})
		if err != nil {
			return nil, err
		}

		return s, nil
	})
}

func (s *TeamServiceServer) Get(ctx context.Context, req *horus.GetTeamRequest) (*horus.Team, error) {
	v, err := s.bare.Team().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	account, err := f.Actor.QueryAccounts().
		Where(account.SiloID(uuid.UUID(v.Silo.Id))).
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	f.ActingAccount = account
	switch account.Role {
	case role.Owner:
		fallthrough
	case role.Admin:
		return v, nil
	}

	_, err = account.QueryMemberships().
		Where(membership.HasTeamWith(team.ID(uuid.UUID(v.Id)))).
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	return v, nil
}

func (s *TeamServiceServer) actorCanModify(ctx context.Context, req *horus.GetTeamRequest) error {
	f := frame.Must(ctx)

	p, err := bare.GetTeamSpecifier(req)
	if err != nil {
		return err
	}

	a, err := f.Actor.QueryAccounts().
		Where(account.HasSiloWith(silo.HasTeamsWith(p))).
		Only(ctx)
	if err != nil {
		return bare.ToStatus(err)
	}
	switch a.Role {
	case role.Owner:
		fallthrough
	case role.Admin:
		return nil
	}

	m, err := a.QueryMemberships().
		Where(membership.HasTeamWith(p)).
		Only(ctx)
	if err != nil {
		return bare.ToStatus(err)
	}
	if m.Role == role.Owner {
		return nil
	}

	return status.Error(codes.PermissionDenied, "team can only be modified by the silo owners, silo admins, or team owners")
}

func (s *TeamServiceServer) Update(ctx context.Context, req *horus.UpdateTeamRequest) (*horus.Team, error) {
	if err := s.actorCanModify(ctx, req.GetKey()); err != nil {
		return nil, err
	}

	return s.bare.Team().Update(ctx, req)
}

func (s *TeamServiceServer) Delete(ctx context.Context, req *horus.GetTeamRequest) (*emptypb.Empty, error) {
	if err := s.actorCanModify(ctx, req); err != nil {
		return nil, err
	}

	return s.bare.Team().Delete(ctx, req)
}

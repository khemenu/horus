package server

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/team"
	"khepri.dev/horus/internal/entutils"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

type TeamServiceServer struct {
	horus.UnimplementedTeamServiceServer
	*base
}

func (s *TeamServiceServer) Create(ctx context.Context, req *horus.CreateTeamRequest) (*horus.Team, error) {
	p, err := s.bare.Silo().Get(ctx, req.GetSilo())
	if err != nil {
		return nil, fmt.Errorf("get silo: %w", err)
	}

	f := frame.Must(ctx)
	v, err := f.Actor.QueryAccounts().
		Where(account.SiloID(uuid.UUID(p.Id))).
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}
	if v.Role != role.Owner {
		return nil, status.Error(codes.PermissionDenied, "only owner can create a team")
	}

	req.Silo = &horus.GetSiloRequest{Key: &horus.GetSiloRequest_Id{Id: p.Id}}
	return entutils.WithTxV(ctx, s.db, func(tx *ent.Tx) (*horus.Team, error) {
		c := tx.Client()
		res, err := bare.NewTeamServiceServer(c).Create(ctx, req)
		if err != nil {
			return nil, err
		}

		_, err = bare.NewMembershipServiceServer(c).Create(ctx, &horus.CreateMembershipRequest{
			Role:    horus.Role_ROLE_OWNER,
			Account: &horus.GetAccountRequest{Id: v.ID[:]},
			Team:    &horus.GetTeamRequest{Id: res.Id},
		})
		if err != nil {
			return nil, err
		}

		return res, nil
	})
}

func (s *TeamServiceServer) Get(ctx context.Context, req *horus.GetTeamRequest) (*horus.Team, error) {
	res, err := s.bare.Team().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	v, err := f.Actor.QueryAccounts().
		Where(account.SiloID(uuid.UUID(res.Silo.Id))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "not found: %s", err)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	f.ActingAccount = v
	if v.Role == role.Owner {
		return res, nil
	}

	_, err = v.QueryMemberships().
		Where(membership.HasTeamWith(team.ID(uuid.UUID(res.Id)))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Errorf(codes.PermissionDenied, "only member who has membership can access: %s", err)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *TeamServiceServer) Update(ctx context.Context, req *horus.UpdateTeamRequest) (*horus.Team, error) {
	v, err := s.Get(ctx, &horus.GetTeamRequest{
		Id: req.GetId(),
	})
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	if f.ActingAccount.Role == role.Owner {
		return s.bare.Team().Update(ctx, req)
	}

	member, err := f.ActingAccount.QueryMemberships().
		Where(membership.HasTeamWith(team.IDEQ(uuid.UUID(v.Id)))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrPermissionDenied
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	if member.Role != role.Owner {
		return nil, ErrPermissionDenied
	}

	return s.bare.Team().Update(ctx, req)
}

func (s *TeamServiceServer) Delete(ctx context.Context, req *horus.DeleteTeamRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "team cannot be deleted manually")
}

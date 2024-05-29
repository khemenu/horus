package server

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/alias"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/team"
	"khepri.dev/horus/internal/entutils"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

type TeamServiceServer struct {
	horus.UnimplementedTeamServiceServer
	*base
}

func (s *TeamServiceServer) Create(ctx context.Context, req *horus.CreateTeamRequest) (*horus.Team, error) {
	silo_id := req.GetTeam().GetSilo().GetId()
	if silo_id == nil {
		return nil, newErrMissingRequiredField("team.silo.id")
	}

	f := frame.Must(ctx)
	v, err := f.Actor.QueryAccounts().
		Where(account.SiloID(uuid.UUID(silo_id))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	if v.Role != account.RoleOWNER {
		return nil, status.Error(codes.PermissionDenied, "only owner can create a team")
	}

	return entutils.WithTxV(ctx, s.db, func(tx *ent.Tx) (*horus.Team, error) {
		c := tx.Client()
		res, err := bare.NewTeamService(c).Create(ctx, &horus.CreateTeamRequest{
			Team: &horus.Team{
				Alias: fx.CoalesceOr(req.Team.GetAlias(), alias.New()),
				Name:  req.Team.GetName(),
				Silo:  req.Team.GetSilo(),

				InterVisibility: horus.Team_INTER_VISIBILITY_PRIVATE,
				IntraVisibility: horus.Team_INTRA_VISIBILITY_PRIVATE,
			},
		})
		if err != nil {
			return nil, err
		}

		_, err = bare.NewMembershipService(c).Create(ctx, &horus.CreateMembershipRequest{
			Membership: &horus.Membership{
				Role:    horus.Membership_ROLE_OWNER,
				Account: &horus.Account{Id: v.ID[:]},
				Team:    &horus.Team{Id: res.Id},
			},
		})
		if err != nil {
			return nil, err
		}

		return res, nil
	})
}

func (s *TeamServiceServer) Get(ctx context.Context, req *horus.GetTeamRequest) (*horus.Team, error) {
	res, err := s.bare.Team().Get(ctx, &horus.GetTeamRequest{
		Id:   req.Id,
		View: horus.GetTeamRequest_WITH_EDGE_IDS,
	})
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
	if v.Role == account.RoleOWNER {
		return res, nil
	}
	if res.InterVisibility == horus.Team_INTER_VISIBILITY_PUBLIC {
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
	team_id := req.Team.GetId()
	if team_id == nil {
		return nil, status.Errorf(codes.InvalidArgument, "required: team.id")
	}

	v, err := s.Get(ctx, &horus.GetTeamRequest{Id: req.GetTeam().GetId()})
	if err != nil {
		return nil, err
	}

	v.Alias = req.Team.Alias
	v.Name = req.Team.Name
	v.Description = req.Team.Description

	f := frame.Must(ctx)
	if f.ActingAccount.Role == account.RoleOWNER {
		v.InterVisibility = req.Team.InterVisibility
		v.IntraVisibility = req.Team.IntraVisibility
		return s.bare.Team().Update(ctx, &horus.UpdateTeamRequest{Team: v})
	}

	member, err := f.ActingAccount.QueryMemberships().
		Where(membership.HasTeamWith(team.ID(uuid.UUID(team_id)))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrPermissionDenied
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	if member.Role == membership.RoleOWNER {
		return s.bare.Team().Update(ctx, &horus.UpdateTeamRequest{Team: v})
	}

	return nil, ErrPermissionDenied
}

func (s *TeamServiceServer) Delete(ctx context.Context, req *horus.DeleteTeamRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "team cannot be deleted manually")
}

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
	f := frame.Must(ctx)

	p, err := bare.GetSiloSpecifier(req.GetSilo())
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
	if v.Role.LowerThan(role.Admin) {
		return nil, status.Error(codes.PermissionDenied, "team can be created only by a silo owner or a silo admin")
	}

	req.Silo = horus.SiloById(v.Edges.Silo.ID)
	return entutils.WithTxV(ctx, s.db, func(tx *ent.Tx) (*horus.Team, error) {
		c := tx.Client()
		res, err := bare.NewTeamServiceServer(c).Create(ctx, req)
		if err != nil {
			return nil, err
		}

		_, err = bare.NewMembershipServiceServer(c).Create(ctx, &horus.CreateMembershipRequest{
			Role:    horus.Role_ROLE_OWNER,
			Account: horus.AccountById(v.ID),
			Team:    horus.TeamByIdV(res.Id),
		})
		if err != nil {
			return nil, err
		}

		return res, nil
	})
}

func (s *TeamServiceServer) Get(ctx context.Context, req *horus.GetTeamRequest) (*horus.Team, error) {
	v, err := s.bare.Team().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	f := frame.Must(ctx)
	acct, err := f.Actor.QueryAccounts().
		Where(account.SiloID(uuid.UUID(v.Silo.Id))).
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	f.ActingAccount = acct
	if acct.Role.HigherThan(role.Member) {
		return v, nil
	}

	_, err = acct.QueryMemberships().
		Where(membership.HasTeamWith(team.ID(uuid.UUID(v.Id)))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Errorf(codes.PermissionDenied, "team can be retrieved only by the team members, silo owners, or silo admins")
		}

		return nil, bare.ToStatus(err)
	}

	return v, nil
}

func (s *TeamServiceServer) Update(ctx context.Context, req *horus.UpdateTeamRequest) (*horus.Team, error) {
	f := frame.Must(ctx)

	v, err := s.Get(ctx, req.GetKey())
	if err != nil {
		return nil, err
	}

	req.Key = horus.TeamByIdV(v.Id)

	account := f.MustGetActingAccount()
	if account.Role.HigherThan(role.Member) {
		// Silo owner and silo admin can update any team in the silo.
		return s.bare.Team().Update(ctx, req)
	}

	membership, err := account.QueryMemberships().
		Where(membership.HasTeamWith(team.IDEQ(uuid.UUID(v.Id)))).
		Only(ctx)
	if err == nil && membership.Role.HigherThan(role.Member) {
		return s.bare.Team().Update(ctx, req)
	}
	if !ent.IsNotFound(err) {
		return nil, bare.ToStatus(err)
	}

	return nil, status.Error(codes.PermissionDenied, "team can be updated only by the owners or admins")
}

func (s *TeamServiceServer) Delete(ctx context.Context, req *horus.GetTeamRequest) (*emptypb.Empty, error) {
	p, err := bare.GetTeamSpecifier(req)
	if err != nil {
		return nil, err
	}

	owners, err := s.db.Membership.Query().
		Where(
			membership.RoleEQ(role.Owner),
			membership.HasTeamWith(p),
		).
		All(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}
	switch len(owners) {
	case 0:
		return nil, status.Errorf(codes.NotFound, "team not found")
	case 1:
		return s.bare.Team().Delete(ctx, req)
	default:
		return nil, status.Error(codes.FailedPrecondition, "only teams with one owner can be deleted.")
	}
}

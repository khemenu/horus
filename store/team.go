package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/log"
	"khepri.dev/horus/store/ent"
	"khepri.dev/horus/store/ent/org"
	"khepri.dev/horus/store/ent/team"
)

func fromEntTeam(v *ent.Team) *horus.Team {
	return &horus.Team{
		Id:    horus.TeamId(v.ID),
		OrgId: horus.OrgId(v.OrgID),
		Name:  v.Name,

		CreatedAt: v.CreatedAt,
	}
}

type teamStore struct {
	*stores
}

func (s *teamStore) New(ctx context.Context, init horus.TeamInit) (*horus.Team, error) {
	return withTx(ctx, s.client, func(tx *ent.Tx) (*horus.Team, error) {
		team, err := tx.Team.Create().
			SetOrgID(uuid.UUID(init.OrgId)).
			SetName(init.Name).
			Save(ctx)
		if err != nil {
			if ent.IsValidationError(err) {
				return nil, errors.Join(horus.ErrInvalidArgument, err)
			}
			if ent.IsConstraintError(err) {
				if strings.Contains(err.Error(), "FOREIGN KEY") {
					return nil, errors.Join(horus.ErrNotExist, err)
				}
				if strings.Contains(err.Error(), "UNIQUE") {
					return nil, errors.Join(horus.ErrExist, err)
				}
			}

			return nil, fmt.Errorf("query: %w", err)
		}

		_, err = s.memberships.new(ctx, tx.Client(), horus.MembershipInit{
			TeamId:   horus.TeamId(team.ID),
			MemberId: init.OwnerId,
			Role:     horus.RoleTeamOwner,
		})
		if err != nil {
			return nil, fmt.Errorf("set team owner: %w", err)
		}

		log.FromCtx(ctx).Info("new team", "id", org.ID)
		return fromEntTeam(team), nil
	})
}

func (s *teamStore) GetById(ctx context.Context, team_id horus.TeamId) (*horus.Team, error) {
	res, err := s.client.Team.Query().
		Where(team.ID(uuid.UUID(team_id))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return fromEntTeam(res), nil
}

func (s *teamStore) GetAllByOrgId(ctx context.Context, org_id horus.OrgId) ([]*horus.Team, error) {
	res, err := s.client.Team.Query().
		Where(team.OrgID(uuid.UUID(org_id))).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return fx.MapV(res, fromEntTeam), nil
}

func (s *teamStore) UpdateById(ctx context.Context, team *horus.Team) (*horus.Team, error) {
	res, err := s.client.Team.UpdateOneID(uuid.UUID(team.Id)).
		SetName(team.Name).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return fromEntTeam(res), nil
}

func (s *teamStore) DeleteByIdFromOrg(ctx context.Context, org_id horus.OrgId, team_id horus.TeamId) error {
	_, err := s.client.Team.Delete().
		Where(team.And(
			team.ID(uuid.UUID(team_id)),
			team.OrgID(uuid.UUID(org_id)),
		)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}

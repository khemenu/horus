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

func team_(v *ent.Team) *horus.Team {
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
			if ent.IsConstraintError(err) {
				if strings.Contains(err.Error(), "FOREIGN KEY") {
					return nil, errors.Join(horus.ErrNotExist, err)
				}
			}

			return nil, fmt.Errorf("save: %w", err)
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
		return team_(team), nil
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

	return team_(res), nil
}

func (s *teamStore) GetAllByOrgId(ctx context.Context, org_id horus.OrgId) ([]*horus.Team, error) {
	res, err := s.client.Team.Query().
		Where(team.OrgID(uuid.UUID(org_id))).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return fx.MapV(res, team_), nil
}

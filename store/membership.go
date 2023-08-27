package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"khepri.dev/horus"
	"khepri.dev/horus/log"
	"khepri.dev/horus/store/ent"
	"khepri.dev/horus/store/ent/membership"
)

func fromEntMembership(v *ent.Membership) *horus.Membership {
	return &horus.Membership{
		TeamId:   horus.TeamId(v.TeamID),
		MemberId: horus.MemberId(v.MemberID),
		Role:     v.Role,

		CreatedAt: v.CreatedAt,
	}
}

type membershipStore struct {
	*stores
}

func (s *membershipStore) new(ctx context.Context, client *ent.Client, init horus.MembershipInit) (*horus.Membership, error) {
	res, err := client.Membership.Create().
		SetTeamID(uuid.UUID(init.TeamId)).
		SetMemberID(uuid.UUID(init.MemberId)).
		SetRole(init.Role).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			if strings.Contains(err.Error(), "FOREIGN KEY") {
				// User or Org does not exist.
				return nil, errors.Join(horus.ErrNotExist, err)
			} else if strings.Contains(err.Error(), "UNIQUE") {
				// Add a user who already a member.
				return nil, errors.Join(horus.ErrExist, err)
			}
		}

		return nil, fmt.Errorf("save: %w", err)
	}

	log.FromCtx(ctx).Info("new membership", "team", init.TeamId, "member", init.MemberId)
	return fromEntMembership(res), nil
}

func (s *membershipStore) New(ctx context.Context, init horus.MembershipInit) (*horus.Membership, error) {
	return s.new(ctx, s.client, init)
}

func (s *membershipStore) GetById(ctx context.Context, team_id horus.TeamId, member_id horus.MemberId) (*horus.Membership, error) {
	res, err := s.client.Membership.Query().
		Where(membership.And(
			membership.TeamID(uuid.UUID(team_id)),
			membership.MemberID(uuid.UUID(member_id)),
		)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return fromEntMembership(res), nil
}

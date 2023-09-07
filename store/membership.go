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
	"khepri.dev/horus/store/ent/member"
	"khepri.dev/horus/store/ent/membership"
	"khepri.dev/horus/store/ent/org"
	"khepri.dev/horus/store/ent/team"
)

func fromEntMembership(v *ent.Membership) *horus.Membership {
	return &horus.Membership{
		TeamId:   horus.TeamId(v.TeamID),
		MemberId: horus.MemberId(v.MemberID),
		Role:     v.Role,

		CreatedAt: v.CreatedAt,
	}
}

func toEntMembership(v *horus.Membership) *ent.Membership {
	return &ent.Membership{
		TeamID:   uuid.UUID(v.TeamId),
		MemberID: uuid.UUID(v.MemberId),
		Role:     v.Role,
	}
}

type membershipStore struct {
	*stores
}

func (s *membershipStore) new(ctx context.Context, client *ent.Client, init horus.MembershipInit) (*horus.Membership, error) {
	return withTx(ctx, s.client, func(tx *ent.Tx) (*horus.Membership, error) {
		ok, err := client.Org.Query().
			Where(org.And(
				org.HasMembersWith(member.ID(uuid.UUID(init.MemberId))),
				org.HasTeamsWith(team.ID(uuid.UUID(init.TeamId))),
			)).
			Exist(ctx)
		if err != nil {
			return nil, fmt.Errorf("query org: %w", err)
		}
		if !ok {
			return nil, horus.ErrNotExist
		}

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
	})
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

func (s *membershipStore) GetByUserIdFromTeam(ctx context.Context, team_id horus.TeamId, user_id horus.UserId) (*horus.Membership, error) {
	res, err := s.client.Membership.Query().
		Where(membership.And(
			membership.TeamID(uuid.UUID(team_id)),
			membership.HasMemberWith(member.UserID(uuid.UUID(user_id))),
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

func (s *membershipStore) GetAllByMemberId(ctx context.Context, member_id horus.MemberId) ([]*horus.Membership, error) {
	res, err := s.client.Membership.Query().
		Where(membership.MemberID(uuid.UUID(member_id))).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return fx.MapV(res, fromEntMembership), nil
}

func (s *membershipStore) UpdateById(ctx context.Context, membership *horus.Membership) (*horus.Membership, error) {
	res, err := s.client.Membership.UpdateOne(toEntMembership(membership)).
		SetRole(membership.Role).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return fromEntMembership(res), nil
}

func (s *membershipStore) DeleteById(ctx context.Context, team_id horus.TeamId, member_id horus.MemberId) error {
	_, err := s.client.Membership.Delete().
		Where(membership.And(
			membership.TeamID(uuid.UUID(team_id)),
			membership.MemberID(uuid.UUID(member_id)),
		)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}

func (s *membershipStore) DeleteByUserIdFromTeam(ctx context.Context, team_id horus.TeamId, user_id horus.UserId) error {
	_, err := s.client.Membership.Delete().
		Where(membership.And(
			membership.TeamID(uuid.UUID(team_id)),
			membership.HasMemberWith(member.UserID(uuid.UUID(user_id))),
		)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}

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
	"khepri.dev/horus/store/ent/identity"
	"khepri.dev/horus/store/ent/member"
	"khepri.dev/horus/store/ent/team"
)

func fromEntMember(v *ent.Member) *horus.Member {
	return &horus.Member{
		Id:     horus.MemberId(v.ID),
		OrgId:  horus.OrgId(v.OrgID),
		UserId: horus.UserId(v.UserID),
		Role:   v.Role,

		Name: v.Name,
		Identities: fx.Associate(v.Edges.Identities, func(v *ent.Identity) (string, *horus.Identity) {
			return v.ID, fromEntIdentity(v)
		}),

		CreatedAt: v.CreatedAt,
	}
}

type memberStore struct {
	*stores
}

func newMember(ctx context.Context, client *ent.Client, init horus.MemberInit) (*horus.Member, error) {
	res, err := client.Member.Create().
		SetOrgID(uuid.UUID(init.OrgId)).
		SetUserID(uuid.UUID(init.UserId)).
		SetRole(init.Role).
		SetName(init.Name).
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

	log.FromCtx(ctx).Info("new user", "id", res.ID)
	return fromEntMember(res), nil
}

func (s *memberStore) New(ctx context.Context, init horus.MemberInit) (*horus.Member, error) {
	return newMember(ctx, s.client, init)
}

func (s *memberStore) GetById(ctx context.Context, member_id horus.MemberId) (*horus.Member, error) {
	res, err := s.client.Member.Query().
		Where(member.ID(uuid.UUID(member_id))).
		WithIdentities().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return fromEntMember(res), nil
}

func (s *memberStore) GetByUserIdFromOrg(ctx context.Context, org_id horus.OrgId, user_id horus.UserId) (*horus.Member, error) {
	res, err := s.client.Member.Query().
		Where(member.And(
			member.OrgID(uuid.UUID(org_id)),
			member.UserID(uuid.UUID(user_id)),
		)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return fromEntMember(res), nil
}

func (s *memberStore) GetByUserIdFromTeam(ctx context.Context, team_id horus.TeamId, user_id horus.UserId) (*horus.Member, error) {
	res, err := s.client.Team.Query().
		Where(team.ID(uuid.UUID(team_id))).
		QueryOrg().
		QueryMembers().
		Where(member.UserID(uuid.UUID(user_id))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return fromEntMember(res), nil
}

func (s *memberStore) GetAllByOrgId(ctx context.Context, org_id horus.OrgId) ([]*horus.Member, error) {
	res, err := s.client.Member.Query().
		Where(member.OrgID(uuid.UUID(org_id))).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	if len(res) == 0 {
		// Org does not exist.
		// There must be at least one member who is owner of the Org.
		return nil, horus.ErrNotExist
	}

	return fx.MapV(res, fromEntMember), nil
}

func (s *memberStore) UpdateById(ctx context.Context, member_ *horus.Member) (*horus.Member, error) {
	return withTx(ctx, s.client, func(tx *ent.Tx) (*horus.Member, error) {
		res, err := tx.Member.UpdateOneID(uuid.UUID(member_.Id)).
			SetName(member_.Name).
			SetRole(member_.Role).
			Save(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, horus.ErrNotExist
			}

			return nil, fmt.Errorf("save: %w", err)
		}

		cnt, err := tx.Member.
			QueryOrg(res).
			QueryMembers().
			Where(member.RoleEQ(horus.RoleOrgOwner)).
			Count(ctx)
		if err != nil {
			return nil, fmt.Errorf("cnt: %w", err)
		}
		if cnt == 0 {
			return nil, fmt.Errorf("%w: organization must have at least one owner", horus.ErrFailedPrecondition)
		}

		return fromEntMember(res), nil
	})
}

func (s *memberStore) AddIdentity(ctx context.Context, member_id horus.MemberId, identity_value horus.IdentityValue) error {
	_, err := withTx(ctx, s.client, func(tx *ent.Tx) (int, error) {
		err := s.client.Member.UpdateOneID(uuid.UUID(member_id)).
			AddIdentityIDs(string(identity_value)).
			Exec(ctx)
		if err != nil {
			if ent.IsConstraintError(err) {
				if strings.Contains(err.Error(), "FOREIGN KEY") {
					return 0, errors.Join(horus.ErrNotExist, err)
				}
			}

			return 0, fmt.Errorf("save: %w", err)
		}

		// Check if same owner.
		cnt, err := s.client.Member.Query().
			Where(member.ID(uuid.UUID(member_id))).
			QueryUser().
			QueryIdentities().
			Where(identity.ID(string(identity_value))).
			Count(ctx)
		if err != nil {
			return 0, fmt.Errorf("query: %w", err)
		}
		if cnt != 1 {
			return 0, fmt.Errorf("%w: different owner", horus.ErrInvalidArgument)
		}

		return 0, nil
	})
	return err
}

func (s *memberStore) AddIdentityByUserIdFromOrg(ctx context.Context, org_id horus.OrgId, user_id horus.UserId, identity_value horus.IdentityValue) error {
	_, err := withTx(ctx, s.client, func(tx *ent.Tx) (int, error) {
		cnt, err := tx.Member.Update().
			Where(member.And(
				member.OrgID(uuid.UUID(org_id)),
				member.UserID(uuid.UUID(user_id)),
			)).
			AddIdentityIDs(string(identity_value)).
			Save(ctx)
		if err != nil {
			if ent.IsConstraintError(err) {
				if strings.Contains(err.Error(), "FOREIGN KEY") {
					return 0, errors.Join(horus.ErrNotExist, err)
				}
			}

			return 0, fmt.Errorf("query member: %w", err)
		}
		if cnt != 1 {
			return 0, horus.ErrNotExist
		}

		cnt, err = tx.Identity.Query().
			Where(identity.And(
				identity.ID(string(identity_value)),
				identity.OwnerID(uuid.UUID(user_id)),
			)).
			Count(ctx)
		if err != nil {
			return 0, fmt.Errorf("query identity: %w", err)
		}
		if cnt != 1 {
			return 0, horus.ErrNotExist
		}

		return 0, nil
	})

	return err
}

func (s *memberStore) RemoveIdentity(ctx context.Context, member_id horus.MemberId, identity_value horus.IdentityValue) error {
	err := s.client.Member.UpdateOneID(uuid.UUID(member_id)).
		RemoveIdentityIDs(string(identity_value)).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil
		}

		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func (s *memberStore) RemoveIdentityByUserIdFromOrg(ctx context.Context, org_id horus.OrgId, user_id horus.UserId, identity_value horus.IdentityValue) error {
	_, err := s.client.Member.Update().
		Where(member.And(
			member.OrgID(uuid.UUID(org_id)),
			member.UserID(uuid.UUID(user_id)),
		)).
		RemoveIdentityIDs(string(identity_value)).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func (s *memberStore) DeleteById(ctx context.Context, member_id horus.MemberId) error {
	err := s.client.Member.DeleteOneID(uuid.UUID(member_id)).Exec(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return fmt.Errorf("save: %w", err)
		}
	}

	return nil
}

func (s *memberStore) DeleteByUserIdFromOrg(ctx context.Context, org_id horus.OrgId, user_id horus.UserId) error {
	_, err := s.client.Member.Delete().
		Where(member.And(
			member.OrgID(uuid.UUID(org_id)),
			member.UserID(uuid.UUID(user_id)),
		)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

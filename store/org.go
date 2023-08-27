package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/log"
	"khepri.dev/horus/store/ent"
	"khepri.dev/horus/store/ent/member"
	"khepri.dev/horus/store/ent/org"
)

func fromEntOrg(v *ent.Org) *horus.Org {
	return &horus.Org{
		Id:   horus.OrgId(v.ID),
		Name: v.Name,

		CreatedAt: v.CreatedAt,
	}
}

type orgStore struct {
	*stores
}

func (s *orgStore) New(ctx context.Context, init horus.OrgInit) (*horus.Org, error) {
	return withTx(ctx, s.client, func(tx *ent.Tx) (*horus.Org, error) {
		org, err := tx.Org.Create().
			SetName(init.Name).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("save: %w", err)
		}

		if _, err := newMember(ctx, tx.Client(), horus.MemberInit{
			OrgId:  horus.OrgId(org.ID),
			UserId: init.OwnerId,
			Role:   horus.RoleOrgOwner,
		}); err != nil {
			return nil, fmt.Errorf("set org owner: %w", err)
		}

		log.FromCtx(ctx).Info("new org", "id", org.ID)
		return fromEntOrg(org), nil
	})
}

func (s *orgStore) GetById(ctx context.Context, org_id horus.OrgId) (*horus.Org, error) {
	res, err := s.client.Org.Query().
		Where(org.ID(uuid.UUID(org_id))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return fromEntOrg(res), nil
}

func (s *orgStore) GetAllByUserId(ctx context.Context, user_id horus.UserId) ([]*horus.Org, error) {
	res, err := s.client.Member.Query().
		Where(member.UserID(uuid.UUID(user_id))).
		QueryOrg().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return fx.MapV(res, fromEntOrg), nil
}

func (s *orgStore) UpdateById(ctx context.Context, org *horus.Org) (*horus.Org, error) {
	res, err := s.client.Org.UpdateOneID(uuid.UUID(org.Id)).
		SetName(org.Name).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("save: %w", err)
	}

	return fromEntOrg(res), nil
}

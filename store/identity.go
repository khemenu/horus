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
	"khepri.dev/horus/store/ent/identity"
)

func identity_(v *ent.Identity) *horus.Identity {
	return &horus.Identity{
		OwnerId: horus.UserId(v.OwnerID),
		Kind:    v.Kind,
		Value:   v.ID,
		Name:    v.Name,

		VerifiedBy: v.VerifiedBy,

		CreatedAt: v.CreatedAt,
	}
}

type identityStore struct {
	*stores
}

func (s *identityStore) new(ctx context.Context, client *ent.Client, init *horus.IdentityInit) (*horus.Identity, error) {
	res, err := client.Identity.Create().
		SetOwnerID(uuid.UUID(init.OwnerId)).
		SetKind(init.Kind).
		SetID(init.Value).
		SetName(init.Name).
		SetVerifiedBy(init.VerifiedBy).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			if strings.Contains(err.Error(), "FOREIGN KEY") {
				return nil, errors.Join(horus.ErrNotExist, err)
			} else {
				return nil, errors.Join(horus.ErrExist, err)
			}
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	log.FromCtx(ctx).Info("new identity")
	return identity_(res), nil
}

func (s *identityStore) New(ctx context.Context, init *horus.IdentityInit) (*horus.Identity, error) {
	if init.OwnerId != horus.UserId(uuid.Nil) {
		return s.new(ctx, s.client, init)
	}

	return withTx(ctx, s.client, func(tx *ent.Tx) (*horus.Identity, error) {
		client := tx.Client()

		owner, err := s.users.new(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("new user: %w", err)
		}

		init.OwnerId = owner.Id
		return s.new(ctx, client, init)
	})
}

func (s *identityStore) GetByValue(ctx context.Context, value string) (*horus.Identity, error) {
	res, err := s.client.Identity.Get(ctx, value)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return identity_(res), nil
}

func (s *identityStore) GetAllByOwner(ctx context.Context, owner_id horus.UserId) (map[string]*horus.Identity, error) {
	res, err := s.client.Identity.Query().
		Where(identity.OwnerID(uuid.UUID(owner_id))).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	rst := map[string]*horus.Identity{}
	for _, v := range res {
		rst[v.ID] = identity_(v)
	}

	return rst, nil
}

func (s *identityStore) Update(ctx context.Context, input *horus.Identity) (*horus.Identity, error) {
	query := s.client.Identity.UpdateOneID(input.Value).
		SetName(input.Name)
	if input.VerifiedBy != "" {
		query.SetVerifiedBy(input.VerifiedBy)
	}
	res, err := query.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return identity_(res), nil
}

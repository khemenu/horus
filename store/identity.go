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

func _identity(identity *ent.Identity) *horus.Identity {
	return &horus.Identity{
		IdentityInit: horus.IdentityInit{
			Value:      identity.ID,
			Name:       identity.Name,
			Kind:       horus.IdentityKind(identity.Kind),
			VerifiedBy: horus.Verifier(identity.VerifiedBy),
		},
		OwnerId:   identity.OwnerID,
		CreatedAt: identity.CreatedAt,
	}
}

type identityStore struct {
	client *ent.Client
}

func NewIdentityStore(client *ent.Client) (horus.IdentityStore, error) {
	s := &identityStore{
		client: client,
	}

	return s, nil
}

func (s *identityStore) Create(ctx context.Context, input *horus.Identity) (*horus.Identity, error) {
	res, err := s.client.Identity.Create().
		SetID(input.Value).
		SetOwnerID(input.OwnerId).
		SetName(input.Name).
		SetKind(string(input.Kind)).
		SetVerifiedBy(string(input.VerifiedBy)).
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
	return _identity(res), nil
}

func (s *identityStore) GetByValue(ctx context.Context, value string) (*horus.Identity, error) {
	res, err := s.client.Identity.Get(ctx, value)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return _identity(res), nil
}

func (s *identityStore) GetAllByOwner(ctx context.Context, owner_id uuid.UUID) (map[string]*horus.Identity, error) {
	res, err := s.client.Identity.Query().
		Where(identity.OwnerID(owner_id)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	rst := map[string]*horus.Identity{}
	for _, v := range res {
		rst[v.ID] = _identity(v)
	}

	return rst, nil
}

func (s *identityStore) Update(ctx context.Context, input *horus.Identity) (*horus.Identity, error) {
	query := s.client.Identity.UpdateOneID(input.Value).
		SetName(input.Name)
	if input.VerifiedBy != "" {
		query.SetVerifiedBy(string(input.VerifiedBy))
	}
	res, err := query.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return _identity(res), nil
}

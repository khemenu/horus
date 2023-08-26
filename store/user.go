package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"khepri.dev/horus"
	"khepri.dev/horus/log"
	"khepri.dev/horus/store/ent"
	"khepri.dev/horus/store/ent/user"
)

func user_(v *ent.User) *horus.User {
	return &horus.User{
		Id:        horus.UserId(v.ID),
		Alias:     v.Alias,
		CreatedAt: v.CreatedAt,
	}
}

type userStore struct {
	client    *ent.Client
	alias_gen horus.Generator
}

type UserStoreOption func(s *userStore) error

func WithCustomUserAlias(generator horus.Generator) UserStoreOption {
	return func(s *userStore) error {
		s.alias_gen = generator
		return nil
	}
}

func NewUserStore(client *ent.Client, opts ...UserStoreOption) (horus.UserStore, error) {
	s := &userStore{
		client:    client,
		alias_gen: horus.DefaultUserAliasGenerator,
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *userStore) New(ctx context.Context) (*horus.User, error) {
	const MaxRetry = 3

	var (
		res *ent.User
		err error
	)
	for i := 0; i < MaxRetry; i++ {
		var alias string
		alias, err = s.alias_gen.New()
		if err != nil {
			return nil, fmt.Errorf("generate alias: %w", err)
		}

		res, err = s.client.User.Create().
			SetAlias(alias).
			Save(ctx)

		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	log.FromCtx(ctx).Info("new user", "id", res.ID)
	return user_(res), nil
}

func (s *userStore) GetById(ctx context.Context, id horus.UserId) (*horus.User, error) {
	res, err := s.client.User.Query().
		Where(user.ID(uuid.UUID(id))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return user_(res), nil
}

func (s *userStore) GetByAlias(ctx context.Context, alias string) (*horus.User, error) {
	res, err := s.client.User.Query().
		Where(user.Alias(alias)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return user_(res), nil
}

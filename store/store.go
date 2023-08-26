package store

import (
	"context"
	"errors"
	"fmt"

	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/store/ent"
)

func NewSqliteMemClient() (*ent.Client, error) {
	return ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
}

type Config struct {
	UserAliasGenerator horus.Generator
	TokenGenerator     horus.Generator
}

type stores struct {
	conf   Config
	client *ent.Client

	users      *userStore
	tokens     *tokenStore
	identities *identityStore

	orgs        *orgStore
	teams       *teamStore
	members     *memberStore
	memberships *membershipStore
}

func (s *stores) Users() horus.UserStore {
	return s.users
}

func (s *stores) Tokens() horus.TokenStore {
	return s.tokens
}

func (s *stores) Identities() horus.IdentityStore {
	return s.identities
}

func (s *stores) Orgs() horus.OrgStore {
	return s.orgs
}

func (s *stores) Teams() horus.TeamStore {
	return s.teams
}

func (s *stores) Members() horus.MemberStore {
	return s.members
}

func (s *stores) Memberships() horus.MembershipStore {
	return s.memberships
}

func NewStores(client *ent.Client, conf *Config) (horus.Stores, error) {
	rst := &stores{
		conf:   *fx.Fallback(conf, &Config{}),
		client: client,
	}

	fx.Default(&rst.conf.UserAliasGenerator, horus.DefaultUserAliasGenerator)
	fx.Default(&rst.conf.TokenGenerator, horus.DefaultOpaqueTokenGenerator)

	rst.users = &userStore{rst}
	rst.tokens = &tokenStore{rst}
	rst.identities = &identityStore{rst}
	rst.orgs = &orgStore{rst}
	rst.teams = &teamStore{rst}
	rst.members = &memberStore{rst}
	rst.memberships = &membershipStore{rst}

	return rst, nil
}

func withTx[R any](ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) (R, error)) (R, error) {
	var rst R

	tx, err := client.Tx(ctx)
	if err != nil {
		return rst, fmt.Errorf("tx begin: %w", err)
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	rst, err = fn(tx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return rst, errors.Join(fmt.Errorf("tx rollback: %w", rerr), err)
		}
		return rst, err
	}
	if err := tx.Commit(); err != nil {
		return rst, fmt.Errorf("tx commit: %w", err)
	}
	return rst, nil
}

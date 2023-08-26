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

type stores struct {
	users      horus.UserStore
	tokens     horus.TokenStore
	identities horus.IdentityStore

	orgs        horus.OrgStore
	teams       horus.TeamStore
	members     horus.MemberStore
	memberships horus.MembershipStore
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

func NewStores(client *ent.Client) (horus.Stores, error) {
	errs := []error{}
	rst := &stores{
		users:       fx.CollectErr(NewUserStore(client)).To(&errs),
		tokens:      fx.CollectErr(NewTokenStore(client)).To(&errs),
		identities:  fx.CollectErr(NewIdentityStore(client)).To(&errs),
		orgs:        fx.CollectErr(NewOrgStore(client)).To(&errs),
		teams:       fx.CollectErr(NewTeamStore(client)).To(&errs),
		members:     fx.CollectErr(NewMemberStore(client)).To(&errs),
		memberships: fx.CollectErr(NewMembershipStore(client)).To(&errs),
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

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

package store

import (
	"errors"

	"khepri.dev/horus"
	"khepri.dev/horus/store/ent"
)

func NewSqliteMemClient() (*ent.Client, error) {
	return ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
}

type stores struct {
	users      horus.UserStore
	tokens     horus.TokenStore
	identities horus.IdentityStore
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

func NewStores(client *ent.Client) (horus.Stores, error) {
	errs := []error{}

	user_store, err := NewUserStore(client)
	if err != nil {
		errs = append(errs, err)
	}

	token_store, err := NewTokenStore(client)
	if err != nil {
		errs = append(errs, err)
	}

	identity_store, err := NewIdentityStore(client)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return &stores{
		users:      user_store,
		tokens:     token_store,
		identities: identity_store,
	}, nil
}

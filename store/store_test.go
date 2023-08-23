package store_test

import (
	"context"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
	"khepri.dev/horus/store"
	"khepri.dev/horus/store/ent"
	"khepri.dev/horus/store/ent/enttest"
)

type storesWrapper struct {
	horus.Stores

	users      horus.UserStore
	tokens     horus.TokenStore
	identities horus.IdentityStore
}

func (s *storesWrapper) Users() horus.UserStore {
	if s.users != nil {
		return s.users
	}

	return s.Stores.Users()
}

func (s *storesWrapper) Tokens() horus.TokenStore {
	if s.tokens != nil {
		return s.tokens
	}

	return s.Stores.Tokens()
}

func (s *storesWrapper) Identities() horus.IdentityStore {
	if s.identities != nil {
		return s.identities
	}

	return s.Stores.Identities()
}

type SuiteWithClient struct {
	suite.Suite
	horus.Stores

	driver_name string
	source_name string
}

func NewSuiteWithClientSqlite() SuiteWithClient {
	return SuiteWithClient{
		driver_name: "sqlite3",
		source_name: "file:ent?mode=memory&cache=shared&_fk=1",
	}
}

func (s *SuiteWithClient) RunWithClient(name string, sub func(require *require.Assertions, ctx context.Context, client *ent.Client)) {
	s.Run(name, func() {
		var (
			require = require.New(s.T())
			ctx     = context.Background()
		)

		client := enttest.Open(s.T(), s.driver_name, s.source_name, enttest.WithOptions(ent.Log(s.T().Log)))
		defer client.Close()

		stores, err := store.NewStores(client)
		require.NoError(err)

		s.Stores = stores
		sub(require, ctx, client)
	})
}

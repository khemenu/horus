package store_test

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
	"khepri.dev/horus/store"
	"khepri.dev/horus/store/ent"
	"khepri.dev/horus/store/ent/enttest"
)

func NewSuiteWithSqliteStores() SuiteWithStores {
	return SuiteWithStores{
		driver_name: "sqlite3",
		source_name: "file:ent?mode=memory&cache=shared&_fk=1",
	}
}

type suiteConfig struct {
	store_conf *store.Config
}

type suiteOption func(opts *suiteConfig)

func withConfig(conf *store.Config) suiteOption {
	return func(opts *suiteConfig) {
		opts.store_conf = conf
	}
}

type SuiteWithStores struct {
	suite.Suite

	driver_name string
	source_name string
}

func (s *SuiteWithStores) RunWithStores(name string, sub func(ctx context.Context, stores horus.Stores), opts ...suiteOption) {
	s.Run(name, func() {
		conf := suiteConfig{}
		for _, opt := range opts {
			opt(&conf)
		}

		client := enttest.Open(s.T(), s.driver_name, s.source_name, enttest.WithOptions(ent.Log(s.T().Log)))
		defer client.Close()

		stores, err := store.NewStores(client, conf.store_conf)
		require.NoError(s.T(), err)

		sub(context.Background(), stores)
	})
}

type SuiteWithStoresUser struct {
	SuiteWithStores

	user *horus.User
}

func (s *SuiteWithStoresUser) RunWithStores(name string, sub func(ctx context.Context, stores horus.Stores), opts ...suiteOption) {
	s.SuiteWithStores.RunWithStores(name, func(ctx context.Context, stores horus.Stores) {
		user, err := stores.Users().New(ctx)
		s.Require().NoError(err)

		s.user = user
		sub(ctx, stores)
	}, opts...)
}

type SuiteWithStoresOrg struct {
	SuiteWithStores

	user  *horus.User
	org   *horus.Org
	owner *horus.Member
}

func (s *SuiteWithStoresOrg) RunWithStores(name string, sub func(ctx context.Context, stores horus.Stores), opts ...suiteOption) {
	s.SuiteWithStores.RunWithStores(name, func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		var err error

		s.user, err = stores.Users().New(ctx)
		require.NoError(err)

		s.org, err = stores.Orgs().New(ctx, horus.OrgInit{OwnerId: s.user.Id})
		require.NoError(err)

		s.owner, err = stores.Members().GetByUserIdFromOrg(ctx, s.org.Id, s.user.Id)
		require.NoError(err)

		sub(ctx, stores)
	}, opts...)
}

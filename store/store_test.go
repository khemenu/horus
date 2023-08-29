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
	horus.Stores

	driver_name string
	source_name string
}

func (s *SuiteWithStores) Run(name string, sub func(ctx context.Context), opts ...suiteOption) {
	s.Suite.Run(name, func() {
		conf := suiteConfig{}
		for _, opt := range opts {
			opt(&conf)
		}

		client := enttest.Open(
			s.T(), s.driver_name, s.source_name,
			enttest.WithOptions(ent.Log(s.T().Log)),
			// enttest.WithOptions(ent.Debug()),
		)
		defer client.Close()

		stores, err := store.NewStores(client, conf.store_conf)
		require.NoError(s.T(), err)

		s.Stores = stores
		sub(context.Background())
	})
}

type SuiteWithStoresUser struct {
	SuiteWithStores

	user *horus.User
}

func (s *SuiteWithStoresUser) InitAmun() *horus.IdentityInit {
	return &horus.IdentityInit{
		OwnerId:    s.user.Id,
		Kind:       horus.IdentityMail,
		Value:      "amun@khepri.dev",
		VerifiedBy: horus.VerifierGoogleOauth2,
	}
}

func (s *SuiteWithStoresUser) InitAtum() *horus.IdentityInit {
	return &horus.IdentityInit{
		OwnerId:    s.user.Id,
		Kind:       horus.IdentityMail,
		Value:      "atum@khepri.dev",
		VerifiedBy: horus.VerifierGoogleOauth2,
	}
}

func (s *SuiteWithStoresUser) Run(name string, sub func(ctx context.Context), opts ...suiteOption) {
	s.SuiteWithStores.Run(name, func(ctx context.Context) {
		user, err := s.Users().New(ctx)
		s.Require().NoError(err)

		s.user = user
		sub(ctx)
	}, opts...)
}

type SuiteWithStoresOrg struct {
	SuiteWithStoresUser

	org   *horus.Org
	owner *horus.Member
}

func (s *SuiteWithStoresOrg) Run(name string, sub func(ctx context.Context), opts ...suiteOption) {
	s.SuiteWithStores.Run(name, func(ctx context.Context) {
		require := s.Require()

		var err error

		s.user, err = s.Users().New(ctx)
		require.NoError(err)

		rst, err := s.Orgs().New(ctx, horus.OrgInit{OwnerId: s.user.Id})
		require.NoError(err)
		s.org = rst.Org

		s.owner, err = s.Members().GetByUserIdFromOrg(ctx, s.org.Id, s.user.Id)
		require.NoError(err)

		sub(ctx)
	}, opts...)
}

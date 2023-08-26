package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
)

type OrgStoreTestSuite struct {
	SuiteWithStores

	user *horus.User
}

func TestOrgStoreSqlite(t *testing.T) {
	suite.Run(t, &OrgStoreTestSuite{
		SuiteWithStores: NewSuiteWithSqliteStores(),
	})
}

func (s *OrgStoreTestSuite) RunWithStores(name string, sub func(ctx context.Context, stores horus.Stores), opts ...suiteOption) {
	s.SuiteWithStores.RunWithStores(name, func(ctx context.Context, stores horus.Stores) {
		user, err := stores.Users().New(ctx)
		s.Require().NoError(err)

		s.user = user
		sub(ctx, stores)
	})
}

func (s *OrgStoreTestSuite) TestNew() {
	s.RunWithStores("with an owner that does not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Orgs().New(ctx, horus.OrgInit{
			OwnerId: horus.UserId(uuid.New()),
			Name:    "khepri",
		})
		require.Error(err, horus.ErrNotExist)
	})

	s.RunWithStores("with an owner that exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		org, err := stores.Orgs().New(ctx, horus.OrgInit{
			OwnerId: s.user.Id,
			Name:    "Khepri",
		})
		require.NoError(err)
		require.Equal("Khepri", org.Name)
	})
}

func (s *OrgStoreTestSuite) TestGetById() {
	s.RunWithStores("not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Orgs().GetById(ctx, horus.OrgId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		expected, err := stores.Orgs().New(ctx, horus.OrgInit{OwnerId: s.user.Id})
		require.NoError(err)

		actual, err := stores.Orgs().GetById(ctx, expected.Id)
		require.NoError(err)
		require.Equal(expected, actual)
	})
}

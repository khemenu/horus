package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
	"khepri.dev/horus/store/ent"
)

type OrgStoreTestSuite struct {
	SuiteWithClient
}

func (suite *OrgStoreTestSuite) RunWithStore(name string, sub func(require *require.Assertions, ctx context.Context, store horus.OrgStore)) {
	suite.RunWithClient(name, func(require *require.Assertions, ctx context.Context, client *ent.Client) {
		sub(require, ctx, suite.Orgs())
	})
}

func TestOrgStoreSqlite(t *testing.T) {
	suite.Run(t, &OrgStoreTestSuite{
		NewSuiteWithClientSqlite(),
	})
}

func (suite *OrgStoreTestSuite) TestNew() {
	suite.RunWithStore("with an owner that does not exist", func(require *require.Assertions, ctx context.Context, store horus.OrgStore) {
		_, err := store.New(ctx, horus.OrgInit{
			OwnerId: horus.UserId(uuid.New()),
			Name:    "khepri",
		})
		require.Error(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with an owner that exists", func(require *require.Assertions, ctx context.Context, store horus.OrgStore) {
		user, err := suite.Users().New(ctx)
		require.NoError(err)

		org, err := store.New(ctx, horus.OrgInit{
			OwnerId: user.Id,
			Name:    "Khepri",
		})
		require.NoError(err)
		require.Equal("Khepri", org.Name)
	})
}

func (suite *OrgStoreTestSuite) TestGetById() {
	suite.RunWithStore("not exist", func(require *require.Assertions, ctx context.Context, store horus.OrgStore) {
		_, err := store.GetById(ctx, horus.OrgId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("exists", func(require *require.Assertions, ctx context.Context, store horus.OrgStore) {
		user, err := suite.Users().New(ctx)
		require.NoError(err)

		expected, err := store.New(ctx, horus.OrgInit{OwnerId: user.Id})
		require.NoError(err)

		actual, err := store.GetById(ctx, expected.Id)
		require.NoError(err)
		require.Equal(expected, actual)
	})
}

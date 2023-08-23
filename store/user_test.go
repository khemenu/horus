package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
	"khepri.dev/horus/store"
	"khepri.dev/horus/store/ent"
)

type UserStoreTestSuite struct {
	SuiteWithClient
}

func (suite *UserStoreTestSuite) RunWithStore(name string, sub func(require *require.Assertions, ctx context.Context, store horus.UserStore), opts ...store.UserStoreOption) {
	suite.RunWithClient(name, func(require *require.Assertions, ctx context.Context, client *ent.Client) {
		if len(opts) > 0 {
			store, err := store.NewUserStore(client, opts...)
			suite.Stores = &storesWrapper{users: store}
			require.NoError(err)
		}

		sub(require, ctx, suite.Users())
	})
}

func TestUserStoreSqlite(t *testing.T) {
	suite.Run(t, &UserStoreTestSuite{
		NewSuiteWithClientSqlite(),
	})
}

func (suite *UserStoreTestSuite) TestNew() {
	suite.RunWithStore("alias collision eventually fail",
		func(require *require.Assertions, ctx context.Context, store horus.UserStore) {
			_, err := store.New(ctx)
			require.NoError(err)

			_, err = store.New(ctx)
			require.Error(err)
		},
		store.WithCustomUserAlias(horus.NewStaticStringGenerator([]rune{'x'}, 1)))
}

func (suite *UserStoreTestSuite) TestGetById() {
	suite.RunWithStore("exists", func(require *require.Assertions, ctx context.Context, store horus.UserStore) {
		expected, err := store.New(ctx)
		require.NoError(err)

		actual, err := store.GetById(ctx, expected.Id)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	suite.RunWithStore("not exist", func(require *require.Assertions, ctx context.Context, store horus.UserStore) {
		_, err := store.GetById(ctx, uuid.New())
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

func (suite *UserStoreTestSuite) TestGetByAlias() {
	suite.RunWithStore("exists", func(require *require.Assertions, ctx context.Context, store horus.UserStore) {
		expected, err := store.New(ctx)
		require.NoError(err)

		actual, err := store.GetByAlias(ctx, expected.Alias)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	suite.RunWithStore("not exist", func(require *require.Assertions, ctx context.Context, store horus.UserStore) {
		_, err := store.GetByAlias(ctx, "not exist")
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

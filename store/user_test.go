package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
	"khepri.dev/horus/store"
)

type UserStoreTestSuite struct {
	SuiteWithStores
}

func TestUserStoreSqlite(t *testing.T) {
	suite.Run(t, &UserStoreTestSuite{
		NewSuiteWithSqliteStores(),
	})
}

func (s *UserStoreTestSuite) TestNew() {
	s.RunWithStores(
		"alias collision eventually fail",
		func(ctx context.Context, stores horus.Stores) {
			require := s.Require()

			_, err := stores.Users().New(ctx)
			require.NoError(err)

			_, err = stores.Users().New(ctx)
			require.Error(err)
		},
		withConfig(&store.Config{
			UserAliasGenerator: horus.NewStaticStringGenerator([]rune{'x'}, 1),
		}),
	)
}

func (s *UserStoreTestSuite) TestGetById() {
	s.RunWithStores("exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		expected, err := stores.Users().New(ctx)
		require.NoError(err)

		actual, err := stores.Users().GetById(ctx, expected.Id)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	s.RunWithStores("not exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Users().GetById(ctx, horus.UserId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

func (s *UserStoreTestSuite) TestGetByAlias() {
	s.RunWithStores("exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		expected, err := stores.Users().New(ctx)
		require.NoError(err)

		actual, err := stores.Users().GetByAlias(ctx, expected.Alias)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	s.RunWithStores("not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Users().GetByAlias(ctx, "not exist")
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

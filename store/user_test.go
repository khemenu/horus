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
	s.Run("alias collision eventually fail",
		func(ctx context.Context) {
			require := s.Require()

			_, err := s.Users().New(ctx)
			require.NoError(err)

			_, err = s.Users().New(ctx)
			require.Error(err)
		},
		withConfig(&store.Config{
			UserAliasGenerator: horus.NewStaticStringGenerator([]rune{'x'}, 1),
		}),
	)
}

func (s *UserStoreTestSuite) TestGetById() {
	s.Run("exists", func(ctx context.Context) {
		require := s.Require()

		expected, err := s.Users().New(ctx)
		require.NoError(err)

		actual, err := s.Users().GetById(ctx, expected.Id)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	s.Run("not exists", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Users().GetById(ctx, horus.UserId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

func (s *UserStoreTestSuite) TestGetByAlias() {
	s.Run("exists", func(ctx context.Context) {
		require := s.Require()

		expected, err := s.Users().New(ctx)
		require.NoError(err)

		actual, err := s.Users().GetByAlias(ctx, expected.Alias)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	s.Run("not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Users().GetByAlias(ctx, "not exist")
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

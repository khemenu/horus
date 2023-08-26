package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
)

type TokenStoreTestSuite struct {
	SuiteWithStoresUser
}

func TestTokenStoreSqlite(t *testing.T) {
	suite.Run(t, &TokenStoreTestSuite{
		SuiteWithStoresUser{
			SuiteWithStores: NewSuiteWithSqliteStores(),
		},
	})
}

func (s *TokenStoreTestSuite) TestGetByValue() {
	s.RunWithStores("exists if the token is not expired", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		expected, err := stores.Tokens().Issue(ctx, horus.TokenInit{
			OwnerId:  s.user.Id,
			Type:     horus.AccessToken,
			Duration: time.Hour,
		})
		require.NoError(err)
		require.Equal(time.Hour, expected.Duration())

		actual, err := stores.Tokens().GetByValue(ctx, expected.Value, horus.AccessToken)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	s.RunWithStores("not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Tokens().GetByValue(ctx, "not exist", horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("not exist if the token is invalid", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Tokens().GetByValue(ctx, "", horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("not exist if the token type is not matched", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		token, err := stores.Tokens().Issue(ctx, horus.TokenInit{
			OwnerId:  s.user.Id,
			Type:     horus.AccessToken,
			Duration: time.Hour,
		})
		require.NoError(err)

		_, err = stores.Tokens().GetByValue(ctx, token.Value, horus.RefreshToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("not exist if the token is expired", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		token, err := stores.Tokens().Issue(ctx, horus.TokenInit{
			OwnerId:  s.user.Id,
			Type:     horus.AccessToken,
			Duration: -time.Hour,
		})
		require.NoError(err)

		_, err = stores.Tokens().GetByValue(ctx, token.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

func (s *TokenStoreTestSuite) TestRevoke() {
	s.RunWithStores("revoked token cannot be get", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		token, err := stores.Tokens().Issue(ctx, horus.TokenInit{
			OwnerId:  s.user.Id,
			Type:     horus.AccessToken,
			Duration: time.Hour,
		})
		require.NoError(err)

		err = stores.Tokens().Revoke(ctx, token.Value)
		require.NoError(err)

		_, err = stores.Tokens().GetByValue(ctx, token.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("not exist not an error", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		err := stores.Tokens().Revoke(ctx, "not exist")
		require.NoError(err)
	})

	s.RunWithStores("invalid value is not an error", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		err := stores.Tokens().Revoke(ctx, "")
		require.NoError(err)
	})
}

func (s *TokenStoreTestSuite) TestRevokeAll() {
	s.RunWithStores("not exist not an error", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		err := stores.Tokens().RevokeAll(ctx, s.user.Id)
		require.NoError(err)
	})

	s.RunWithStores("all tokens are revoked", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		init := horus.TokenInit{
			OwnerId:  s.user.Id,
			Type:     horus.RefreshToken,
			Duration: time.Hour,
		}
		foo, err := stores.Tokens().Issue(ctx, init)
		require.NoError(err)

		bar, err := stores.Tokens().Issue(ctx, init)
		require.NoError(err)

		init.Type = horus.AccessToken
		baz, err := stores.Tokens().Issue(ctx, init)
		require.NoError(err)

		qux, err := stores.Tokens().Issue(ctx, init)
		require.NoError(err)

		err = stores.Tokens().RevokeAll(ctx, s.user.Id)
		require.NoError(err)

		_, err = stores.Tokens().GetByValue(ctx, foo.Value, horus.RefreshToken)
		require.ErrorIs(err, horus.ErrNotExist)

		_, err = stores.Tokens().GetByValue(ctx, bar.Value, horus.RefreshToken)
		require.ErrorIs(err, horus.ErrNotExist)

		_, err = stores.Tokens().GetByValue(ctx, baz.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)

		_, err = stores.Tokens().GetByValue(ctx, qux.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

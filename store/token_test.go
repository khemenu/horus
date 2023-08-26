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
	s.Run("exists if the token is not expired", func(ctx context.Context) {
		require := s.Require()

		expected, err := s.Tokens().Issue(ctx, horus.TokenInit{
			OwnerId:  s.user.Id,
			Type:     horus.AccessToken,
			Duration: time.Hour,
		})
		require.NoError(err)
		require.Equal(time.Hour, expected.Duration())

		actual, err := s.Tokens().GetByValue(ctx, expected.Value, horus.AccessToken)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	s.Run("not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Tokens().GetByValue(ctx, "not exist", horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("not exist if the token is invalid", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Tokens().GetByValue(ctx, "", horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("not exist if the token type is not matched", func(ctx context.Context) {
		require := s.Require()

		token, err := s.Tokens().Issue(ctx, horus.TokenInit{
			OwnerId:  s.user.Id,
			Type:     horus.AccessToken,
			Duration: time.Hour,
		})
		require.NoError(err)

		_, err = s.Tokens().GetByValue(ctx, token.Value, horus.RefreshToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("not exist if the token is expired", func(ctx context.Context) {
		require := s.Require()

		token, err := s.Tokens().Issue(ctx, horus.TokenInit{
			OwnerId:  s.user.Id,
			Type:     horus.AccessToken,
			Duration: -time.Hour,
		})
		require.NoError(err)

		_, err = s.Tokens().GetByValue(ctx, token.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

func (s *TokenStoreTestSuite) TestRevoke() {
	s.Run("revoked token cannot be get", func(ctx context.Context) {
		require := s.Require()

		token, err := s.Tokens().Issue(ctx, horus.TokenInit{
			OwnerId:  s.user.Id,
			Type:     horus.AccessToken,
			Duration: time.Hour,
		})
		require.NoError(err)

		err = s.Tokens().Revoke(ctx, token.Value)
		require.NoError(err)

		_, err = s.Tokens().GetByValue(ctx, token.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("not exist not an error", func(ctx context.Context) {
		require := s.Require()

		err := s.Tokens().Revoke(ctx, "not exist")
		require.NoError(err)
	})

	s.Run("invalid value is not an error", func(ctx context.Context) {
		require := s.Require()

		err := s.Tokens().Revoke(ctx, "")
		require.NoError(err)
	})
}

func (s *TokenStoreTestSuite) TestRevokeAll() {
	s.Run("not exist not an error", func(ctx context.Context) {
		require := s.Require()

		err := s.Tokens().RevokeAll(ctx, s.user.Id)
		require.NoError(err)
	})

	s.Run("all tokens are revoked", func(ctx context.Context) {
		require := s.Require()

		init := horus.TokenInit{
			OwnerId:  s.user.Id,
			Type:     horus.RefreshToken,
			Duration: time.Hour,
		}
		foo, err := s.Tokens().Issue(ctx, init)
		require.NoError(err)

		bar, err := s.Tokens().Issue(ctx, init)
		require.NoError(err)

		init.Type = horus.AccessToken
		baz, err := s.Tokens().Issue(ctx, init)
		require.NoError(err)

		qux, err := s.Tokens().Issue(ctx, init)
		require.NoError(err)

		err = s.Tokens().RevokeAll(ctx, s.user.Id)
		require.NoError(err)

		_, err = s.Tokens().GetByValue(ctx, foo.Value, horus.RefreshToken)
		require.ErrorIs(err, horus.ErrNotExist)

		_, err = s.Tokens().GetByValue(ctx, bar.Value, horus.RefreshToken)
		require.ErrorIs(err, horus.ErrNotExist)

		_, err = s.Tokens().GetByValue(ctx, baz.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)

		_, err = s.Tokens().GetByValue(ctx, qux.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

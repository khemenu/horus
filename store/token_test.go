package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
	"khepri.dev/horus/store"
	"khepri.dev/horus/store/ent"
)

type TokenStoreTestSuite struct {
	SuiteWithClient
}

func (suite *TokenStoreTestSuite) RunWithStore(name string, sub func(require *require.Assertions, ctx context.Context, store horus.TokenStore), opts ...store.TokenStoreOption) {
	suite.RunWithClient(name, func(require *require.Assertions, ctx context.Context, client *ent.Client) {
		if len(opts) > 0 {
			store, err := store.NewTokenStore(client, opts...)
			suite.Stores = &storesWrapper{tokens: store}
			require.NoError(err)
		}

		sub(require, ctx, suite.Tokens())
	})
}

func TestTokenStoreSqlite(t *testing.T) {
	suite.Run(t, &TokenStoreTestSuite{
		SuiteWithClient: NewSuiteWithClientSqlite(),
	})
}

func (suite *TokenStoreTestSuite) TestGetByValue() {
	suite.RunWithStore("exists if the token is not expired", func(require *require.Assertions, ctx context.Context, store horus.TokenStore) {
		user, err := suite.Users().New(ctx)
		require.NoError(err)

		expected, err := store.Issue(ctx, horus.TokenInit{
			OwnerId:  user.Id,
			Type:     horus.AccessToken,
			Duration: time.Hour,
		})
		require.NoError(err)
		require.Equal(time.Hour, expected.Duration())

		actual, err := store.GetByValue(ctx, expected.Value, horus.AccessToken)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	suite.RunWithStore("not exist", func(require *require.Assertions, ctx context.Context, store horus.TokenStore) {
		_, err := store.GetByValue(ctx, "not exist", horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("not exist if the token is invalid", func(require *require.Assertions, ctx context.Context, store horus.TokenStore) {
		_, err := store.GetByValue(ctx, "", horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("not exist if the token type is not matched", func(require *require.Assertions, ctx context.Context, store horus.TokenStore) {
		user, err := suite.Users().New(ctx)
		require.NoError(err)

		token, err := store.Issue(ctx, horus.TokenInit{
			OwnerId:  user.Id,
			Type:     horus.AccessToken,
			Duration: time.Hour,
		})
		require.NoError(err)

		_, err = store.GetByValue(ctx, token.Value, horus.RefreshToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("not exist if the token is expired", func(require *require.Assertions, ctx context.Context, store horus.TokenStore) {
		user, err := suite.Users().New(ctx)
		require.NoError(err)

		token, err := store.Issue(ctx, horus.TokenInit{
			OwnerId:  user.Id,
			Type:     horus.AccessToken,
			Duration: -time.Hour,
		})
		require.NoError(err)

		_, err = store.GetByValue(ctx, token.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

func (suite *TokenStoreTestSuite) TestRevoke() {
	suite.RunWithStore("revoked token cannot be get", func(require *require.Assertions, ctx context.Context, store horus.TokenStore) {
		user, err := suite.Users().New(ctx)
		require.NoError(err)

		token, err := store.Issue(ctx, horus.TokenInit{
			OwnerId:  user.Id,
			Type:     horus.AccessToken,
			Duration: time.Hour,
		})
		require.NoError(err)

		err = store.Revoke(ctx, token.Value)
		require.NoError(err)

		_, err = store.GetByValue(ctx, token.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("not exist not an error", func(require *require.Assertions, ctx context.Context, store horus.TokenStore) {
		err := store.Revoke(ctx, "not exist")
		require.NoError(err)
	})

	suite.RunWithStore("invalid value is not an error", func(require *require.Assertions, ctx context.Context, store horus.TokenStore) {
		err := store.Revoke(ctx, "")
		require.NoError(err)
	})
}

func (suite *TokenStoreTestSuite) TestRevokeAll() {
	suite.RunWithStore("not exist not an error", func(require *require.Assertions, ctx context.Context, store horus.TokenStore) {
		user, err := suite.Users().New(ctx)
		require.NoError(err)

		err = store.RevokeAll(ctx, uuid.UUID(user.Id))
		require.NoError(err)
	})

	suite.RunWithStore("all tokens are revoked", func(require *require.Assertions, ctx context.Context, store horus.TokenStore) {
		user, err := suite.Users().New(ctx)
		require.NoError(err)

		init := horus.TokenInit{
			OwnerId:  user.Id,
			Type:     horus.RefreshToken,
			Duration: time.Hour,
		}
		foo, err := suite.Tokens().Issue(ctx, init)
		require.NoError(err)

		bar, err := suite.Tokens().Issue(ctx, init)
		require.NoError(err)

		init.Type = horus.AccessToken
		baz, err := suite.Tokens().Issue(ctx, init)
		require.NoError(err)

		qux, err := suite.Tokens().Issue(ctx, init)
		require.NoError(err)

		err = store.RevokeAll(ctx, uuid.UUID(user.Id))
		require.NoError(err)

		_, err = store.GetByValue(ctx, foo.Value, horus.RefreshToken)
		require.ErrorIs(err, horus.ErrNotExist)

		_, err = store.GetByValue(ctx, bar.Value, horus.RefreshToken)
		require.ErrorIs(err, horus.ErrNotExist)

		_, err = store.GetByValue(ctx, baz.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)

		_, err = store.GetByValue(ctx, qux.Value, horus.AccessToken)
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

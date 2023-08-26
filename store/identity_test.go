package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
	"khepri.dev/horus/store/ent"
)

type IdentityStoreTestSuite struct {
	SuiteWithClient
}

var (
	init_ra = horus.IdentityInit{
		Value:      "ra@khepri.dev",
		Kind:       horus.IdentityEmail,
		VerifiedBy: horus.VerifierGoogleOauth2,
	}
	init_atum = horus.IdentityInit{
		Value:      "atum@khepri.dev",
		Kind:       horus.IdentityEmail,
		VerifiedBy: horus.VerifierGoogleOauth2,
	}
)

func (suite *IdentityStoreTestSuite) RunWithStore(name string, sub func(require *require.Assertions, ctx context.Context, store horus.IdentityStore)) {
	suite.RunWithClient(name, func(require *require.Assertions, ctx context.Context, client *ent.Client) {
		sub(require, ctx, suite.Identities())
	})
}

func (suite *IdentityStoreTestSuite) CreateUser(ctx context.Context) *horus.User {
	rst, err := suite.Users().New(ctx)
	require.NoError(suite.T(), err)

	return rst
}

func (suite *IdentityStoreTestSuite) CreateIdentity(ctx context.Context, init horus.IdentityInit) *horus.Identity {
	user := suite.CreateUser(ctx)

	rst, err := suite.Identities().Create(ctx, &horus.Identity{
		IdentityInit: init,
		OwnerId:      user.Id,
	})
	require.NoError(suite.T(), err)

	return rst
}

func TestIdentityStoreSqlite(t *testing.T) {
	suite.Run(t, &IdentityStoreTestSuite{
		SuiteWithClient: NewSuiteWithClientSqlite(),
	})
}

func (suite *IdentityStoreTestSuite) TestCreate() {
	suite.RunWithStore("user can have multiple identity", func(require *require.Assertions, ctx context.Context, store horus.IdentityStore) {
		user := suite.CreateUser(ctx)

		_, err := store.Create(ctx, &horus.Identity{
			IdentityInit: init_ra,
			OwnerId:      user.Id,
		})
		require.NoError(err)

		_, err = store.Create(ctx, &horus.Identity{
			IdentityInit: init_atum,
			OwnerId:      user.Id,
		})
		require.NoError(err)
	})

	suite.RunWithStore("owner must be exist", func(require *require.Assertions, ctx context.Context, store horus.IdentityStore) {
		_, err := store.Create(ctx, &horus.Identity{
			IdentityInit: init_ra,
			OwnerId:      horus.UserId(uuid.New()),
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("value is unique across users", func(require *require.Assertions, ctx context.Context, store horus.IdentityStore) {
		user1 := suite.CreateUser(ctx)
		user2 := suite.CreateUser(ctx)

		_, err := store.Create(ctx, &horus.Identity{
			IdentityInit: init_ra,
			OwnerId:      user1.Id,
		})
		require.NoError(err)

		_, err = store.Create(ctx, &horus.Identity{
			IdentityInit: init_ra,
			OwnerId:      user2.Id,
		})
		require.ErrorIs(err, horus.ErrExist)
	})
}

func (suite *IdentityStoreTestSuite) TestGetByValue() {
	suite.RunWithStore("exists", func(require *require.Assertions, ctx context.Context, store horus.IdentityStore) {
		expected := suite.CreateIdentity(ctx, init_ra)

		actual, err := store.GetByValue(ctx, expected.Value)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	suite.RunWithStore("not exist", func(require *require.Assertions, ctx context.Context, store horus.IdentityStore) {
		_, err := store.GetByValue(ctx, "not exist")
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

func (suite *IdentityStoreTestSuite) TestGetAllByOwner() {
	suite.RunWithStore("exists", func(require *require.Assertions, ctx context.Context, store horus.IdentityStore) {
		user := suite.CreateUser(ctx)

		ra, err := store.Create(ctx, &horus.Identity{
			IdentityInit: init_ra,
			OwnerId:      user.Id,
		})
		require.NoError(err)

		atum, err := store.Create(ctx, &horus.Identity{
			IdentityInit: init_atum,
			OwnerId:      user.Id,
		})
		require.NoError(err)

		identities, err := store.GetAllByOwner(ctx, user.Id)
		require.NoError(err)

		require.Equal(map[string]*horus.Identity{
			ra.Value:   ra,
			atum.Value: atum,
		}, identities)
	})

	suite.RunWithStore("not exist", func(require *require.Assertions, ctx context.Context, store horus.IdentityStore) {
		identities, err := store.GetAllByOwner(ctx, horus.UserId(uuid.New()))
		require.NoError(err)
		require.Empty(identities)
	})
}

func (suite *IdentityStoreTestSuite) TestUpdate() {
	suite.RunWithStore("exists", func(require *require.Assertions, ctx context.Context, store horus.IdentityStore) {
		expected := suite.CreateIdentity(ctx, init_ra)

		expected.VerifiedBy = "something else"
		actual, err := store.Update(ctx, expected)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	suite.RunWithStore("not exist", func(require *require.Assertions, ctx context.Context, store horus.IdentityStore) {
		_, err := store.Update(ctx, &horus.Identity{
			IdentityInit: horus.IdentityInit{
				Value: "not exist",
			},
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

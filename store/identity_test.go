package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
)

type IdentityStoreTestSuite struct {
	SuiteWithStoresUser
}

func TestIdentityStoreSqlite(t *testing.T) {
	suite.Run(t, &IdentityStoreTestSuite{
		SuiteWithStoresUser{
			SuiteWithStores: NewSuiteWithSqliteStores(),
		},
	})
}

func (s *IdentityStoreTestSuite) InitAmun() *horus.IdentityInit {
	return &horus.IdentityInit{
		OwnerId:    s.user.Id,
		Kind:       horus.IdentityEmail,
		Value:      "amun@khepri.dev",
		VerifiedBy: horus.VerifierGoogleOauth2,
	}
}

func (s *IdentityStoreTestSuite) InitAtum() *horus.IdentityInit {
	return &horus.IdentityInit{
		OwnerId:    s.user.Id,
		Kind:       horus.IdentityEmail,
		Value:      "atum@khepri.dev",
		VerifiedBy: horus.VerifierGoogleOauth2,
	}
}

func (s *IdentityStoreTestSuite) TestNew() {
	s.RunWithStores("user can have multiple identity", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		_, err = stores.Identities().New(ctx, s.InitAtum())
		require.NoError(err)
	})

	s.RunWithStores("owner must be exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		init := s.InitAmun()
		init.OwnerId = horus.UserId(uuid.New())

		_, err := stores.Identities().New(ctx, init)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("value must be unique across users", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		init := s.InitAmun()
		_, err := stores.Identities().New(ctx, init)
		require.NoError(err)

		other, err := stores.Users().New(ctx)
		require.NoError(err)

		init.OwnerId = other.Id
		_, err = stores.Identities().New(ctx, init)
		require.ErrorIs(err, horus.ErrExist)
	})

	s.RunWithStores("new user is created if owner ID not given", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		init := s.InitAmun()
		init.OwnerId = horus.UserId(uuid.Nil)
		identity, err := stores.Identities().New(ctx, init)
		require.NoError(err)

		_, err = stores.Users().GetById(ctx, identity.OwnerId)
		require.NoError(err)
	})
}

func (s *IdentityStoreTestSuite) TestGetByValue() {
	s.RunWithStores("exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		expected, err := stores.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		actual, err := stores.Identities().GetByValue(ctx, expected.Value)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	s.RunWithStores("not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Identities().GetByValue(ctx, "not exist")
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

func (s *IdentityStoreTestSuite) TestGetAllByOwner() {
	s.RunWithStores("exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		amun, err := stores.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		atum, err := stores.Identities().New(ctx, s.InitAtum())
		require.NoError(err)

		identities, err := stores.Identities().GetAllByOwner(ctx, s.user.Id)
		require.NoError(err)
		require.Equal(map[string]*horus.Identity{
			amun.Value: amun,
			atum.Value: atum,
		}, identities)
	})

	s.RunWithStores("not exist not an error", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		identities, err := stores.Identities().GetAllByOwner(ctx, horus.UserId(uuid.New()))
		require.NoError(err)
		require.Empty(identities)
	})
}

func (s *IdentityStoreTestSuite) TestUpdate() {
	s.RunWithStores("exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		expected, err := stores.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		expected.VerifiedBy = "something else"
		actual, err := stores.Identities().Update(ctx, expected)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	s.RunWithStores("not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Identities().Update(ctx, &horus.Identity{Value: "not exist"})
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

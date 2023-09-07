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

func (s *IdentityStoreTestSuite) TestNew() {
	s.Run("user can have multiple identity", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		_, err = s.Identities().New(ctx, s.InitAtum())
		require.NoError(err)
	})

	s.Run("owner must be exist", func(ctx context.Context) {
		require := s.Require()

		init := s.InitAmun()
		init.OwnerId = horus.UserId(uuid.New())

		_, err := s.Identities().New(ctx, init)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("value must be unique across users", func(ctx context.Context) {
		require := s.Require()

		init := s.InitAmun()
		_, err := s.Identities().New(ctx, init)
		require.NoError(err)

		other, err := s.Users().New(ctx)
		require.NoError(err)

		init.OwnerId = other.Id
		_, err = s.Identities().New(ctx, init)
		require.ErrorIs(err, horus.ErrExist)
	})

	s.Run("new user is created if owner ID not given", func(ctx context.Context) {
		require := s.Require()

		init := s.InitAmun()
		init.OwnerId = horus.UserId(uuid.Nil)
		identity, err := s.Identities().New(ctx, init)
		require.NoError(err)

		_, err = s.Users().GetById(ctx, identity.OwnerId)
		require.NoError(err)
	})
}

func (s *IdentityStoreTestSuite) TestGetByValue() {
	s.Run("exists", func(ctx context.Context) {
		require := s.Require()

		expected, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		actual, err := s.Identities().GetByValue(ctx, expected.Value)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	s.Run("not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Identities().GetByValue(ctx, "not exist")
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

func (s *IdentityStoreTestSuite) TestGetAllByOwner() {
	s.Run("exists", func(ctx context.Context) {
		require := s.Require()

		amun, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		atum, err := s.Identities().New(ctx, s.InitAtum())
		require.NoError(err)

		identities, err := s.Identities().GetAllByOwner(ctx, s.user.Id)
		require.NoError(err)
		require.Equal(map[horus.IdentityValue]*horus.Identity{
			amun.Value: amun,
			atum.Value: atum,
		}, identities)
	})

	s.Run("not exist not an error", func(ctx context.Context) {
		require := s.Require()

		identities, err := s.Identities().GetAllByOwner(ctx, horus.UserId(uuid.New()))
		require.NoError(err)
		require.Empty(identities)
	})
}

func (s *IdentityStoreTestSuite) TestUpdate() {
	s.Run("exists", func(ctx context.Context) {
		require := s.Require()

		expected, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		expected.VerifiedBy = "something else"
		actual, err := s.Identities().Update(ctx, expected)
		require.NoError(err)
		require.Equal(expected, actual)
	})

	s.Run("not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Identities().Update(ctx, &horus.Identity{Value: "not exist"})
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

func (s *IdentityStoreTestSuite) TestDelete() {
	s.Run("exists", func(ctx context.Context) {
		require := s.Require()

		amun, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		err = s.Identities().Delete(ctx, amun.Value)
		require.NoError(err)

		_, err = s.Identities().GetByValue(ctx, amun.Value)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("not exist", func(ctx context.Context) {
		require := s.Require()

		err := s.Identities().Delete(ctx, "not exist")
		require.NoError(err)
	})

	s.Run("identity of member also deleted", func(ctx context.Context) {
		require := s.Require()

		amun, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		rst, err := s.Orgs().New(ctx, horus.OrgInit{
			OwnerId: s.user.Id,
			Name:    "Khepri",
		})
		require.NoError(err)

		err = s.Members().AddIdentity(ctx, rst.Owner.Id, amun.Value)
		require.NoError(err)

		err = s.Identities().Delete(ctx, amun.Value)
		require.NoError(err)

		member, err := s.Members().GetById(ctx, rst.Owner.Id)
		require.NoError(err)
		require.Empty(member.Identities)
	})
}

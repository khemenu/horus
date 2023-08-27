package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
)

type OrgStoreTestSuite struct {
	SuiteWithStoresUser
}

func TestOrgStoreSqlite(t *testing.T) {
	suite.Run(t, &OrgStoreTestSuite{
		SuiteWithStoresUser{
			SuiteWithStores: NewSuiteWithSqliteStores(),
		},
	})
}

func (s *OrgStoreTestSuite) TestNew() {
	s.Run("with an owner that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Orgs().New(ctx, horus.OrgInit{
			OwnerId: horus.UserId(uuid.New()),
			Name:    "khepri",
		})
		require.Error(err, horus.ErrNotExist)
	})

	s.Run("with an owner that exists", func(ctx context.Context) {
		require := s.Require()

		org, err := s.Orgs().New(ctx, horus.OrgInit{
			OwnerId: s.user.Id,
			Name:    "Khepri",
		})
		require.NoError(err)
		require.Equal("Khepri", org.Name)
	})
}

func (s *OrgStoreTestSuite) TestGetById() {
	s.Run("not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Orgs().GetById(ctx, horus.OrgId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("exists", func(ctx context.Context) {
		require := s.Require()

		expected, err := s.Orgs().New(ctx, horus.OrgInit{OwnerId: s.user.Id})
		require.NoError(err)

		actual, err := s.Orgs().GetById(ctx, expected.Id)
		require.NoError(err)
		require.Equal(expected, actual)
	})
}

func (s *OrgStoreTestSuite) TestGetAllByUserId() {
	s.Run("user does not exist", func(ctx context.Context) {
		require := s.Require()

		orgs, err := s.Orgs().GetAllByUserId(ctx, horus.UserId(uuid.New()))
		require.NoError(err)
		require.Empty(orgs)
	})

	s.Run("user does not belongs to any orgs", func(ctx context.Context) {
		require := s.Require()

		orgs, err := s.Orgs().GetAllByUserId(ctx, s.user.Id)
		require.NoError(err)
		require.Empty(orgs)
	})

	s.Run("user belongs to many orgs", func(ctx context.Context) {
		require := s.Require()

		org1, err := s.Orgs().New(ctx, horus.OrgInit{OwnerId: s.user.Id})
		require.NoError(err)

		org2, err := s.Orgs().New(ctx, horus.OrgInit{OwnerId: s.user.Id})
		require.NoError(err)

		orgs, err := s.Orgs().GetAllByUserId(ctx, s.user.Id)
		require.NoError(err)
		require.ElementsMatch([]*horus.Org{org1, org2}, orgs)
	})
}

func (s *OrgStoreTestSuite) TestUpdateById() {
	s.Run("org does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Orgs().UpdateById(ctx, &horus.Org{
			Id:   horus.OrgId(uuid.New()),
			Name: "foo",
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("org exists", func(ctx context.Context) {
		require := s.Require()

		org, err := s.Orgs().New(ctx, horus.OrgInit{OwnerId: s.user.Id})
		require.NoError(err)
		require.Empty(org.Name)

		updated, err := s.Orgs().UpdateById(ctx, &horus.Org{
			Id:   org.Id,
			Name: "foo",
		})
		require.NoError(err)
		require.Equal("foo", updated.Name)
	})
}

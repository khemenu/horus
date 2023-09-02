package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
)

type TeamStoreTestSuite struct {
	SuiteWithStoresOrg
}

func TestTeamStoreSqlite(t *testing.T) {
	suite.Run(t, &TeamStoreTestSuite{
		SuiteWithStoresOrg: SuiteWithStoresOrg{
			SuiteWithStoresUser: SuiteWithStoresUser{
				SuiteWithStores: NewSuiteWithSqliteStores(),
			},
		},
	})
}

func (s *TeamStoreTestSuite) TestNew() {
	s.Run("with an org that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Teams().New(ctx, horus.TeamInit{
			OrgId:   horus.OrgId(uuid.New()),
			OwnerId: s.owner.Id,
			Name:    "A",
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("with a user that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: horus.MemberId(uuid.New()),
			Name:    "A",
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("with an org and a user that exists", func(ctx context.Context) {
		require := s.Require()

		team, err := s.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: s.owner.Id,
			Name:    "A",
		})
		require.NoError(err)
		require.Equal(s.org.Id, team.OrgId)
		require.Equal("A", team.Name)

		membership, err := s.Memberships().GetById(ctx, team.Id, s.owner.Id)
		require.NoError(err)
		require.Equal(team.Id, membership.TeamId)
		require.Equal(s.owner.Id, membership.MemberId)
		require.Equal(horus.RoleTeamOwner, membership.Role)
	})
}

func (s *TeamStoreTestSuite) TestGetById() {
	s.Run("not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Teams().GetById(ctx, horus.TeamId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("exists", func(ctx context.Context) {
		require := s.Require()

		expected, err := s.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: s.owner.Id,
			Name:    "A",
		})
		require.NoError(err)

		actual, err := s.Teams().GetById(ctx, expected.Id)
		require.NoError(err)
		require.Equal(expected, actual)
	})
}

func (s *TeamStoreTestSuite) TestGetAllByOrgId() {
	s.Run("not exist", func(ctx context.Context) {
		require := s.Require()

		teams, err := s.Teams().GetAllByOrgId(ctx, horus.OrgId(uuid.New()))
		require.NoError(err)
		require.Empty(teams)
	})

	s.Run("exists", func(ctx context.Context) {
		require := s.Require()

		a, err := s.Teams().New(ctx, horus.TeamInit{OrgId: s.org.Id, OwnerId: s.owner.Id, Name: "A"})
		require.NoError(err)

		b, err := s.Teams().New(ctx, horus.TeamInit{OrgId: s.org.Id, OwnerId: s.owner.Id, Name: "B"})
		require.NoError(err)

		teams, err := s.Teams().GetAllByOrgId(ctx, s.org.Id)
		require.NoError(err)
		require.ElementsMatch([]*horus.Team{a, b}, teams)
	})
}

func (s *TeamStoreTestSuite) TestUpdateById() {
	s.Run("team does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Teams().UpdateById(ctx, &horus.Team{
			Id:    horus.TeamId(uuid.New()),
			OrgId: s.owner.OrgId,
			Name:  "A-team",
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("team exists", func(ctx context.Context) {
		require := s.Require()

		team, err := s.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.owner.OrgId,
			OwnerId: s.owner.Id,
			Name:    "A-team",
		})
		require.NoError(err)

		team.Name = "Crazy 88s"
		team, err = s.Teams().UpdateById(ctx, team)
		require.NoError(err)
		require.Equal("Crazy 88s", team.Name)

		actual, err := s.Teams().GetById(ctx, team.Id)
		require.NoError(err)
		require.Equal(team, actual)
	})
}

func (s *TeamStoreTestSuite) TestDeleteByIdFromOrg() {
	s.Run("team does not exist", func(ctx context.Context) {
		require := s.Require()

		err := s.Teams().DeleteByIdFromOrg(ctx, s.org.Id, horus.TeamId(uuid.New()))
		require.NoError(err)
	})

	s.Run("team does not exist in given org", func(ctx context.Context) {
		require := s.Require()

		rst, err := s.Orgs().New(ctx, horus.OrgInit{OwnerId: s.user.Id})
		require.NoError(err)

		team, err := s.Teams().New(ctx, horus.TeamInit{
			OrgId:   rst.Org.Id,
			OwnerId: rst.Owner.Id,
		})
		require.NoError(err)

		err = s.Teams().DeleteByIdFromOrg(ctx, s.org.Id, team.Id)
		require.NoError(err)

		_, err = s.Teams().GetById(ctx, team.Id)
		require.NoError(err)
	})

	s.Run("team exists", func(ctx context.Context) {
		require := s.Require()

		team, err := s.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: s.owner.Id,
		})
		require.NoError(err)

		err = s.Teams().DeleteByIdFromOrg(ctx, s.org.Id, team.Id)
		require.NoError(err)

		_, err = s.Teams().GetById(ctx, team.Id)
		require.ErrorIs(err, horus.ErrNotExist)
	})
}

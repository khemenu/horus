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
			SuiteWithStores: NewSuiteWithSqliteStores(),
		},
	})
}

func (s *TeamStoreTestSuite) TestNew() {
	s.RunWithStores("with an org that does not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Teams().New(ctx, horus.TeamInit{
			OrgId:   horus.OrgId(uuid.New()),
			OwnerId: s.owner.Id,
			Name:    "A",
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("with a user that does not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: horus.MemberId(uuid.New()),
			Name:    "A",
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("with an org and a user that exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		team, err := stores.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: s.owner.Id,
			Name:    "A",
		})
		require.NoError(err)
		require.Equal(s.org.Id, team.OrgId)
		require.Equal("A", team.Name)

		membership, err := stores.Memberships().GetById(ctx, team.Id, s.owner.Id)
		require.NoError(err)
		require.Equal(team.Id, membership.TeamId)
		require.Equal(s.owner.Id, membership.MemberId)
		require.Equal(horus.RoleTeamOwner, membership.Role)
	})
}

func (s *TeamStoreTestSuite) TestGetById() {
	s.RunWithStores("not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Teams().GetById(ctx, horus.TeamId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		expected, err := stores.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: s.owner.Id,
			Name:    "A",
		})
		require.NoError(err)

		actual, err := stores.Teams().GetById(ctx, expected.Id)
		require.NoError(err)
		require.Equal(expected, actual)
	})
}

func (s *TeamStoreTestSuite) TestGetAllByOrgId() {
	s.RunWithStores("not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		teams, err := stores.Teams().GetAllByOrgId(ctx, horus.OrgId(uuid.New()))
		require.NoError(err)
		require.Empty(teams)
	})

	s.RunWithStores("exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		a, err := stores.Teams().New(ctx, horus.TeamInit{OrgId: s.org.Id, OwnerId: s.owner.Id, Name: "A"})
		require.NoError(err)

		b, err := stores.Teams().New(ctx, horus.TeamInit{OrgId: s.org.Id, OwnerId: s.owner.Id, Name: "B"})
		require.NoError(err)

		teams, err := stores.Teams().GetAllByOrgId(ctx, s.org.Id)
		require.NoError(err)
		require.ElementsMatch([]*horus.Team{a, b}, teams)
	})
}

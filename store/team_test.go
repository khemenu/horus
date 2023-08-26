package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
	"khepri.dev/horus/store/ent"
)

type TeamStoreTestSuite struct {
	SuiteWithClient

	user  *horus.User
	org   *horus.Org
	owner *horus.Member
}

func (suite *TeamStoreTestSuite) RunWithStore(name string, sub func(require *require.Assertions, ctx context.Context, store horus.TeamStore)) {
	suite.RunWithClient(name, func(require *require.Assertions, ctx context.Context, client *ent.Client) {
		var err error
		suite.user, err = suite.Users().New(ctx)
		require.NoError(err)

		suite.org, err = suite.Orgs().New(ctx, horus.OrgInit{OwnerId: suite.user.Id})
		require.NoError(err)

		suite.owner, err = suite.Members().GetByUserIdFromOrg(ctx, suite.org.Id, suite.user.Id)
		require.NoError(err)

		sub(require, ctx, suite.Teams())
	})
}

func TestTeamStoreSqlite(t *testing.T) {
	suite.Run(t, &TeamStoreTestSuite{
		SuiteWithClient: NewSuiteWithClientSqlite(),
	})
}

func (suite *TeamStoreTestSuite) TestNew() {
	suite.RunWithStore("with an org that does not exist", func(require *require.Assertions, ctx context.Context, store horus.TeamStore) {
		_, err := store.New(ctx, horus.TeamInit{
			OrgId:   horus.OrgId(uuid.New()),
			OwnerId: suite.owner.Id,
			Name:    "A",
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with a user that does not exist", func(require *require.Assertions, ctx context.Context, store horus.TeamStore) {
		_, err := store.New(ctx, horus.TeamInit{
			OrgId:   suite.org.Id,
			OwnerId: horus.MemberId(uuid.New()),
			Name:    "A",
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with an org and a user that exists", func(require *require.Assertions, ctx context.Context, store horus.TeamStore) {
		team, err := store.New(ctx, horus.TeamInit{
			OrgId:   suite.org.Id,
			OwnerId: suite.owner.Id,
			Name:    "A",
		})
		require.NoError(err)
		require.Equal(suite.org.Id, team.OrgId)
		require.Equal("A", team.Name)

		membership, err := suite.Memberships().GetById(ctx, team.Id, suite.owner.Id)
		require.NoError(err)
		require.Equal(team.Id, membership.TeamId)
		require.Equal(suite.owner.Id, membership.MemberId)
		require.Equal(horus.RoleTeamOwner, membership.Role)
	})
}

func (suite *TeamStoreTestSuite) TestGetById() {
	suite.RunWithStore("not exist", func(require *require.Assertions, ctx context.Context, store horus.TeamStore) {
		_, err := store.GetById(ctx, horus.TeamId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("exists", func(require *require.Assertions, ctx context.Context, store horus.TeamStore) {
		expected, err := store.New(ctx, horus.TeamInit{
			OrgId:   suite.org.Id,
			OwnerId: suite.owner.Id,
			Name:    "A",
		})
		require.NoError(err)

		actual, err := store.GetById(ctx, expected.Id)
		require.NoError(err)
		require.Equal(expected, actual)
	})
}

func (suite *TeamStoreTestSuite) TestGetAllByOrgId() {
	suite.RunWithStore("not exist", func(require *require.Assertions, ctx context.Context, store horus.TeamStore) {
		teams, err := store.GetAllByOrgId(ctx, horus.OrgId(uuid.New()))
		require.NoError(err)
		require.Empty(teams)
	})

	suite.RunWithStore("exists", func(require *require.Assertions, ctx context.Context, store horus.TeamStore) {
		a, err := store.New(ctx, horus.TeamInit{OrgId: suite.org.Id, OwnerId: suite.owner.Id, Name: "A"})
		require.NoError(err)

		b, err := store.New(ctx, horus.TeamInit{OrgId: suite.org.Id, OwnerId: suite.owner.Id, Name: "B"})
		require.NoError(err)

		teams, err := store.GetAllByOrgId(ctx, suite.org.Id)
		require.NoError(err)
		require.ElementsMatch([]*horus.Team{a, b}, teams)
	})
}

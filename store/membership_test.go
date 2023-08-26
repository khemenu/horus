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

type MembershipStoreTestSuite struct {
	SuiteWithClient

	user  *horus.User
	org   *horus.Org
	owner *horus.Member
}

func (suite *MembershipStoreTestSuite) RunWithStore(name string, sub func(require *require.Assertions, ctx context.Context, store horus.MembershipStore)) {
	suite.RunWithClient(name, func(require *require.Assertions, ctx context.Context, client *ent.Client) {
		var err error
		suite.user, err = suite.Users().New(ctx)
		require.NoError(err)

		suite.org, err = suite.Orgs().New(ctx, horus.OrgInit{OwnerId: suite.user.Id})
		require.NoError(err)

		suite.owner, err = suite.Members().GetByUserIdFromOrg(ctx, suite.org.Id, suite.user.Id)
		require.NoError(err)

		sub(require, ctx, suite.Memberships())
	})
}

func TestMembershipStoreSqlite(t *testing.T) {
	suite.Run(t, &MembershipStoreTestSuite{
		SuiteWithClient: NewSuiteWithClientSqlite(),
	})
}

func (suite *MembershipStoreTestSuite) TestNew() {
	suite.RunWithStore("with a team that does not exist", func(require *require.Assertions, ctx context.Context, store horus.MembershipStore) {
		_, err := store.New(ctx, horus.MembershipInit{
			TeamId:   horus.TeamId(uuid.New()),
			MemberId: horus.MemberId(suite.user.Id),
			Role:     horus.RoleTeamMember,
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with a member that does not exist", func(require *require.Assertions, ctx context.Context, store horus.MembershipStore) {
		team, err := suite.Teams().New(ctx, horus.TeamInit{
			OrgId:   suite.org.Id,
			OwnerId: suite.owner.Id,
		})
		require.NoError(err)

		_, err = store.New(ctx, horus.MembershipInit{
			TeamId:   team.Id,
			MemberId: horus.MemberId(uuid.New()),
			Role:     horus.RoleTeamMember,
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with a member that already in a membership", func(require *require.Assertions, ctx context.Context, store horus.MembershipStore) {
		team, err := suite.Teams().New(ctx, horus.TeamInit{
			OrgId:   suite.org.Id,
			OwnerId: suite.owner.Id,
		})
		require.NoError(err)

		_, err = store.New(ctx, horus.MembershipInit{
			TeamId:   team.Id,
			MemberId: suite.owner.Id,
			Role:     horus.RoleTeamMember,
		})
		require.ErrorIs(err, horus.ErrExist)
	})

	suite.RunWithStore("with a team and a member that exists", func(require *require.Assertions, ctx context.Context, store horus.MembershipStore) {
		user2, err := suite.Users().New(ctx)
		require.NoError(err)

		member, err := suite.Members().New(ctx, horus.MemberInit{
			OrgId:  suite.org.Id,
			UserId: user2.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		team, err := suite.Teams().New(ctx, horus.TeamInit{
			OrgId:   suite.org.Id,
			OwnerId: suite.owner.Id,
		})
		require.NoError(err)

		membership, err := store.New(ctx, horus.MembershipInit{
			TeamId:   team.Id,
			MemberId: member.Id,
			Role:     horus.RoleTeamMember,
		})
		require.NoError(err)
		require.Equal(horus.RoleTeamMember, membership.Role)
	})
}

func (suite *MembershipStoreTestSuite) TestGetById() {
	suite.RunWithStore("with a team that does not exist", func(require *require.Assertions, ctx context.Context, store horus.MembershipStore) {
		_, err := store.GetById(ctx, horus.TeamId(uuid.New()), suite.owner.Id)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with a member that does not exist", func(require *require.Assertions, ctx context.Context, store horus.MembershipStore) {
		team, err := suite.Teams().New(ctx, horus.TeamInit{
			OrgId:   suite.org.Id,
			OwnerId: suite.owner.Id,
		})
		require.NoError(err)

		_, err = store.GetById(ctx, team.Id, horus.MemberId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with a team and a member that exists", func(require *require.Assertions, ctx context.Context, store horus.MembershipStore) {
		team, err := suite.Teams().New(ctx, horus.TeamInit{
			OrgId:   suite.org.Id,
			OwnerId: suite.owner.Id,
		})
		require.NoError(err)

		membership, err := store.GetById(ctx, team.Id, suite.owner.Id)
		require.NoError(err)
		require.Equal(horus.RoleTeamOwner, membership.Role)
	})
}

// func (suite *MembershipStoreTestSuite) TestGetAllByOrgId() {
// 	suite.RunWithStore("not exist", func(require *require.Assertions, ctx context.Context, store horus.MembershipStore) {
// 		teams, err := store.GetAllByOrgId(ctx, uuid.New())
// 		require.NoError(err)
// 		require.Empty(teams)
// 	})

// 	suite.RunWithStore("exists", func(require *require.Assertions, ctx context.Context, store horus.MembershipStore) {
// 		user, err := suite.Users().New(ctx)
// 		require.NoError(err)

// 		org, err := suite.Orgs().New(ctx, horus.OrgInit{OwnerId: user.Id})
// 		require.NoError(err)

// 		member, err := suite.Members().GetByUserId(ctx, user.Id)
// 		require.NoError(err)

// 		a, err := store.New(ctx, horus.MembershipInit{OrgId: org.Id, OwnerId: member.Id, Name: "A"})
// 		require.NoError(err)

// 		b, err := store.New(ctx, horus.MembershipInit{OrgId: org.Id, OwnerId: member.Id, Name: "B"})
// 		require.NoError(err)

// 		teams, err := store.GetAllByOrgId(ctx, org.Id)
// 		require.NoError(err)
// 		require.ElementsMatch([]*horus.Membership{a, b}, teams)
// 	})
// }

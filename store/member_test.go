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

type MemberStoreTestSuite struct {
	SuiteWithClient

	user *horus.User
	org  *horus.Org
}

func (suite *MemberStoreTestSuite) RunWithStore(name string, sub func(require *require.Assertions, ctx context.Context, store horus.MemberStore)) {
	suite.RunWithClient(name, func(require *require.Assertions, ctx context.Context, client *ent.Client) {
		var err error
		suite.user, err = suite.Users().New(ctx)
		require.NoError(err)

		suite.org, err = suite.Orgs().New(ctx, horus.OrgInit{OwnerId: suite.user.Id})
		require.NoError(err)

		sub(require, ctx, suite.Members())
	})
}

func TestMemberStoreSqlite(t *testing.T) {
	suite.Run(t, &MemberStoreTestSuite{
		SuiteWithClient: NewSuiteWithClientSqlite(),
	})
}

func (suite *MemberStoreTestSuite) TestNew() {
	suite.RunWithStore("with an org that does not exist", func(require *require.Assertions, ctx context.Context, store horus.MemberStore) {
		user, err := suite.Users().New(ctx)
		require.NoError(err)

		_, err = store.New(ctx, horus.MemberInit{
			OrgId:  horus.OrgId(uuid.New()),
			UserId: user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.Error(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with a user that does not exist", func(require *require.Assertions, ctx context.Context, store horus.MemberStore) {
		_, err := store.New(ctx, horus.MemberInit{
			OrgId:  suite.org.Id,
			UserId: horus.UserId(uuid.New()),
			Role:   horus.RoleOrgMember,
		})
		require.Error(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with a user that already a member", func(require *require.Assertions, ctx context.Context, store horus.MemberStore) {
		_, err := store.New(ctx, horus.MemberInit{
			OrgId:  suite.org.Id,
			UserId: suite.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.ErrorIs(err, horus.ErrExist)
	})

	suite.RunWithStore("with an org and a user that exists", func(require *require.Assertions, ctx context.Context, store horus.MemberStore) {
		user, err := suite.Users().New(ctx)
		require.NoError(err)

		member, err := store.New(ctx, horus.MemberInit{
			OrgId:  suite.org.Id,
			UserId: user.Id,
			Role:   horus.RoleOrgMember,
			Name:   "Khepri",
		})
		require.NoError(err)
		require.Equal(horus.RoleOrgMember, member.Role)
		require.Equal("Khepri", member.Name)
	})
}

func (suite *MemberStoreTestSuite) TestGetByUserIdFromOrg() {
	suite.RunWithStore("from an org that does not exist", func(require *require.Assertions, ctx context.Context, store horus.MemberStore) {
		_, err := store.GetByUserIdFromOrg(ctx, horus.OrgId(uuid.New()), suite.user.Id)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with a member that does not exist", func(require *require.Assertions, ctx context.Context, store horus.MemberStore) {
		_, err := store.GetByUserIdFromOrg(ctx, suite.org.Id, horus.UserId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("from an org with member that exists", func(require *require.Assertions, ctx context.Context, store horus.MemberStore) {
		member, err := store.GetByUserIdFromOrg(ctx, suite.org.Id, suite.user.Id)
		require.NoError(err)
		require.Equal(member.Role, horus.RoleOrgOwner)

		user, err := suite.Users().GetById(ctx, member.UserId)
		require.NoError(err)
		require.Equal(suite.user, user)
	})
}

func (suite *MemberStoreTestSuite) TestGetAllByOrgId() {
	suite.RunWithStore("with org that does not exist", func(require *require.Assertions, ctx context.Context, store horus.MemberStore) {
		_, err := store.GetAllByOrgId(ctx, horus.OrgId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	suite.RunWithStore("with org that exists", func(require *require.Assertions, ctx context.Context, store horus.MemberStore) {
		owner, err := suite.Members().GetByUserIdFromOrg(ctx, suite.org.Id, suite.user.Id)
		require.NoError(err)

		user, err := suite.Users().New(ctx)
		require.NoError(err)

		member, err := suite.Members().New(ctx, horus.MemberInit{
			OrgId:  suite.org.Id,
			UserId: user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		members, err := store.GetAllByOrgId(ctx, suite.org.Id)
		require.NoError(err)
		require.ElementsMatch([]*horus.Member{owner, member}, members)
	})
}

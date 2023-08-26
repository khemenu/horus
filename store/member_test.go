package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
)

type MemberStoreTestSuite struct {
	SuiteWithStoresOrg
}

func TestMemberStoreSqlite(t *testing.T) {
	suite.Run(t, &MemberStoreTestSuite{
		SuiteWithStoresOrg: SuiteWithStoresOrg{
			SuiteWithStores: NewSuiteWithSqliteStores(),
		},
	})
}

func (s *MemberStoreTestSuite) TestNew() {
	s.Run("with an org that does not exist", func(ctx context.Context) {
		require := s.Require()

		user, err := s.Users().New(ctx)
		require.NoError(err)

		_, err = s.Members().New(ctx, horus.MemberInit{
			OrgId:  horus.OrgId(uuid.New()),
			UserId: user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.Error(err, horus.ErrNotExist)
	})

	s.Run("with a user that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Members().New(ctx, horus.MemberInit{
			OrgId:  s.org.Id,
			UserId: horus.UserId(uuid.New()),
			Role:   horus.RoleOrgMember,
		})
		require.Error(err, horus.ErrNotExist)
	})

	s.Run("with a user that already a member", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Members().New(ctx, horus.MemberInit{
			OrgId:  s.org.Id,
			UserId: s.user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.ErrorIs(err, horus.ErrExist)
	})

	s.Run("with an org and a user that exists", func(ctx context.Context) {
		require := s.Require()

		user, err := s.Users().New(ctx)
		require.NoError(err)

		member, err := s.Members().New(ctx, horus.MemberInit{
			OrgId:  s.org.Id,
			UserId: user.Id,
			Role:   horus.RoleOrgMember,
			Name:   "Khepri",
		})
		require.NoError(err)
		require.Equal(horus.RoleOrgMember, member.Role)
		require.Equal("Khepri", member.Name)
	})
}

func (s *MemberStoreTestSuite) TestGetByUserIdFromOrg() {
	s.Run("from an org that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Members().GetByUserIdFromOrg(ctx, horus.OrgId(uuid.New()), s.user.Id)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("with a member that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Members().GetByUserIdFromOrg(ctx, s.org.Id, horus.UserId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("from an org with member that exists", func(ctx context.Context) {
		require := s.Require()

		member, err := s.Members().GetByUserIdFromOrg(ctx, s.org.Id, s.user.Id)
		require.NoError(err)
		require.Equal(member.Role, horus.RoleOrgOwner)

		user, err := s.Users().GetById(ctx, member.UserId)
		require.NoError(err)
		require.Equal(s.user, user)
	})
}

func (s *MemberStoreTestSuite) TestGetAllByOrgId() {
	s.Run("with org that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Members().GetAllByOrgId(ctx, horus.OrgId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("with org that exists", func(ctx context.Context) {
		require := s.Require()

		owner, err := s.Members().GetByUserIdFromOrg(ctx, s.org.Id, s.user.Id)
		require.NoError(err)

		user, err := s.Users().New(ctx)
		require.NoError(err)

		member, err := s.Members().New(ctx, horus.MemberInit{
			OrgId:  s.org.Id,
			UserId: user.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		members, err := s.Members().GetAllByOrgId(ctx, s.org.Id)
		require.NoError(err)
		require.ElementsMatch([]*horus.Member{owner, member}, members)
	})
}

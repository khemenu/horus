package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
)

type MembershipStoreTestSuite struct {
	SuiteWithStoresOrg

	team *horus.Team
}

func (s *MembershipStoreTestSuite) RunWithStores(name string, sub func(ctx context.Context, stores horus.Stores), opts ...suiteOption) {
	s.SuiteWithStoresOrg.RunWithStores(name, func(ctx context.Context, stores horus.Stores) {
		team, err := stores.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: s.owner.Id,
		})
		s.Require().NoError(err)

		s.team = team
		sub(ctx, stores)
	})
}

func TestMembershipStoreSqlite(t *testing.T) {
	suite.Run(t, &MembershipStoreTestSuite{
		SuiteWithStoresOrg: SuiteWithStoresOrg{
			SuiteWithStores: NewSuiteWithSqliteStores(),
		},
	})
}

func (s *MembershipStoreTestSuite) TestNew() {
	s.RunWithStores("with a team that does not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Memberships().New(ctx, horus.MembershipInit{
			TeamId:   horus.TeamId(uuid.New()),
			MemberId: s.owner.Id,
			Role:     horus.RoleTeamMember,
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("with a member that does not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Memberships().New(ctx, horus.MembershipInit{
			TeamId:   s.team.Id,
			MemberId: horus.MemberId(uuid.New()),
			Role:     horus.RoleTeamMember,
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("with a member that already in a membership", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Memberships().New(ctx, horus.MembershipInit{
			TeamId:   s.team.Id,
			MemberId: s.owner.Id,
			Role:     horus.RoleTeamMember,
		})
		require.ErrorIs(err, horus.ErrExist)
	})

	s.RunWithStores("with a team and a member that exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		user2, err := stores.Users().New(ctx)
		require.NoError(err)

		member, err := stores.Members().New(ctx, horus.MemberInit{
			OrgId:  s.org.Id,
			UserId: user2.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		membership, err := stores.Memberships().New(ctx, horus.MembershipInit{
			TeamId:   s.team.Id,
			MemberId: member.Id,
			Role:     horus.RoleTeamMember,
		})
		require.NoError(err)
		require.Equal(horus.RoleTeamMember, membership.Role)
	})
}

func (s *MembershipStoreTestSuite) TestGetById() {
	s.RunWithStores("with a team that does not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Memberships().GetById(ctx, horus.TeamId(uuid.New()), s.owner.Id)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("with a member that does not exist", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		_, err := stores.Memberships().GetById(ctx, s.team.Id, horus.MemberId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.RunWithStores("with a team and a member that exists", func(ctx context.Context, stores horus.Stores) {
		require := s.Require()

		team, err := stores.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: s.owner.Id,
		})
		require.NoError(err)

		membership, err := stores.Memberships().GetById(ctx, team.Id, s.owner.Id)
		require.NoError(err)
		require.Equal(horus.RoleTeamOwner, membership.Role)
	})
}

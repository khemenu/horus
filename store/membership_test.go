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

func (s *MembershipStoreTestSuite) Run(name string, sub func(ctx context.Context), opts ...suiteOption) {
	s.SuiteWithStoresOrg.Run(name, func(ctx context.Context) {
		team, err := s.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: s.owner.Id,
		})
		s.Require().NoError(err)

		s.team = team
		sub(ctx)
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
	s.Run("with a team that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Memberships().New(ctx, horus.MembershipInit{
			TeamId:   horus.TeamId(uuid.New()),
			MemberId: s.owner.Id,
			Role:     horus.RoleTeamMember,
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("with a member that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Memberships().New(ctx, horus.MembershipInit{
			TeamId:   s.team.Id,
			MemberId: horus.MemberId(uuid.New()),
			Role:     horus.RoleTeamMember,
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("with a member that already in a membership", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Memberships().New(ctx, horus.MembershipInit{
			TeamId:   s.team.Id,
			MemberId: s.owner.Id,
			Role:     horus.RoleTeamMember,
		})
		require.ErrorIs(err, horus.ErrExist)
	})

	s.Run("with a team and a member that exists", func(ctx context.Context) {
		require := s.Require()

		user2, err := s.Users().New(ctx)
		require.NoError(err)

		member, err := s.Members().New(ctx, horus.MemberInit{
			OrgId:  s.org.Id,
			UserId: user2.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		membership, err := s.Memberships().New(ctx, horus.MembershipInit{
			TeamId:   s.team.Id,
			MemberId: member.Id,
			Role:     horus.RoleTeamMember,
		})
		require.NoError(err)
		require.Equal(horus.RoleTeamMember, membership.Role)
	})
}

func (s *MembershipStoreTestSuite) TestGetById() {
	s.Run("with a team that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Memberships().GetById(ctx, horus.TeamId(uuid.New()), s.owner.Id)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("with a member that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Memberships().GetById(ctx, s.team.Id, horus.MemberId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("with a team and a member that exists", func(ctx context.Context) {
		require := s.Require()

		team, err := s.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: s.owner.Id,
		})
		require.NoError(err)

		membership, err := s.Memberships().GetById(ctx, team.Id, s.owner.Id)
		require.NoError(err)
		require.Equal(horus.RoleTeamOwner, membership.Role)
	})
}

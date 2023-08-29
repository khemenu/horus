package store_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
)

type MemberStoreTestSuite struct {
	SuiteWithStoresOrg
}

func TestMemberStoreSqlite(t *testing.T) {
	suite.Run(t, &MemberStoreTestSuite{
		SuiteWithStoresOrg: SuiteWithStoresOrg{
			SuiteWithStoresUser: SuiteWithStoresUser{
				SuiteWithStores: NewSuiteWithSqliteStores(),
			},
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

func (s *MemberStoreTestSuite) TestGetById() {
	s.Run("member that does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Members().GetById(ctx, horus.MemberId(uuid.New()))
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("member that exists", func(ctx context.Context) {
		require := s.Require()

		member, err := s.Members().GetById(ctx, s.owner.Id)
		require.NoError(err)
		require.Equal(member.Role, horus.RoleOrgOwner)
	})

	s.Run("with identities", func(ctx context.Context) {
		require := s.Require()

		amun, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		atum, err := s.Identities().New(ctx, s.InitAtum())
		require.NoError(err)

		err = s.Members().AddIdentity(ctx, s.owner.Id, amun.Value)
		require.NoError(err)

		err = s.Members().AddIdentity(ctx, s.owner.Id, atum.Value)
		require.NoError(err)

		member, err := s.Members().GetById(ctx, s.owner.Id)
		require.NoError(err)
		require.Equal(
			map[string]*horus.Identity{
				string(amun.Value): amun,
				string(atum.Value): atum,
			},
			member.Identities,
		)
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

func (s *MemberStoreTestSuite) TestUpdateById() {
	s.Run("member does not exist", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Members().UpdateById(ctx, &horus.Member{
			Id:   horus.MemberId(uuid.New()),
			Role: horus.RoleOrgMember,
		})
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("member exists", func(ctx context.Context) {
		require := s.Require()

		updated, err := s.Members().UpdateById(ctx, &horus.Member{
			Id:   s.owner.Id,
			Role: horus.RoleOrgOwner,
			Name: "foo",
		})
		require.NoError(err)
		require.Equal("foo", updated.Name)
	})

	s.Run("sole owner to member", func(ctx context.Context) {
		require := s.Require()

		_, err := s.Members().UpdateById(ctx, &horus.Member{
			Id:   s.owner.Id,
			Role: horus.RoleOrgMember,
		})
		require.ErrorIs(err, horus.ErrFailedPrecondition)

		owner, err := s.Members().GetById(ctx, s.owner.Id)
		require.NoError(err)
		require.Equal(horus.RoleOrgOwner, owner.Role)
	})
}

func (s *MemberStoreTestSuite) TestAddIdentity() {
	s.Run("member does not exist", func(ctx context.Context) {
		require := s.Require()

		identity, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		err = s.Members().AddIdentity(ctx, horus.MemberId(uuid.New()), identity.Value)
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("identity does not exist", func(ctx context.Context) {
		require := s.Require()

		err := s.Members().AddIdentity(ctx, s.owner.Id, "not exists")
		require.ErrorIs(err, horus.ErrNotExist)
	})

	s.Run("member and identity both exists", func(ctx context.Context) {
		require := s.Require()

		identity, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		err = s.Members().AddIdentity(ctx, s.owner.Id, identity.Value)
		require.NoError(err)
	})

	s.Run("already exist", func(ctx context.Context) {
		require := s.Require()

		identity, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		err = s.Members().AddIdentity(ctx, s.owner.Id, identity.Value)
		require.NoError(err)

		err = s.Members().AddIdentity(ctx, s.owner.Id, identity.Value)
		require.NoError(err)
	})

	s.Run("different owner", func(ctx context.Context) {
		require := s.Require()

		other, err := s.Users().New(ctx)
		require.NoError(err)

		other_member, err := s.Members().New(ctx, horus.MemberInit{
			OrgId:  s.org.Id,
			UserId: other.Id,
			Role:   horus.RoleOrgMember,
		})
		require.NoError(err)

		identity, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		err = s.Members().AddIdentity(ctx, other_member.Id, identity.Value)
		require.ErrorIs(err, horus.ErrInvalidArgument)
	})
}

func (s *MemberStoreTestSuite) TestRemoveIdentity() {
	for _, tc := range []struct {
		desc           string
		member_id      horus.MemberId
		identity_value horus.IdentityValue
	}{
		{
			desc:      "member does not exist",
			member_id: horus.MemberId(uuid.New()),
		},
		{
			desc:           "identity does not exist",
			identity_value: "not exists",
		},
		{
			desc:           "member and identity both does not exist",
			member_id:      horus.MemberId(uuid.New()),
			identity_value: "not exists",
		},
	} {
		s.Run(tc.desc, func(ctx context.Context) {
			require := s.Require()

			identity, err := s.Identities().New(ctx, s.InitAmun())
			require.NoError(err)

			err = s.Members().AddIdentity(ctx, s.owner.Id, identity.Value)
			require.NoError(err)

			member_id := fx.Fallback(tc.member_id, s.owner.Id)
			identity_value := fx.Fallback(tc.identity_value, identity.Value)

			err = s.Members().RemoveIdentity(ctx, member_id, identity_value)
			require.NoError(err)
		})
	}

	s.Run("member and identity both exists", func(ctx context.Context) {
		require := s.Require()

		identity, err := s.Identities().New(ctx, s.InitAmun())
		require.NoError(err)

		err = s.Members().AddIdentity(ctx, s.owner.Id, identity.Value)
		require.NoError(err)

		err = s.Members().RemoveIdentity(ctx, s.owner.Id, identity.Value)
		require.NoError(err)

		member, err := s.Members().GetById(ctx, s.owner.Id)
		require.NoError(err)
		require.NotContains(member.Identities, identity.Value)
	})
}

func (s *MemberStoreTestSuite) TestDeleteById() {
	s.Run("member does not exist", func(ctx context.Context) {
		require := s.Require()

		err := s.Members().DeleteById(ctx, horus.MemberId(uuid.New()))
		require.NoError(err)
	})

	s.Run("member exists", func(ctx context.Context) {
		require := s.Require()

		err := s.Members().DeleteById(ctx, s.owner.Id)
		require.NoError(err)
	})

	s.Run("membership is also deleted", func(ctx context.Context) {
		require := s.Require()

		team, err := s.Teams().New(ctx, horus.TeamInit{
			OrgId:   s.org.Id,
			OwnerId: s.owner.Id,
		})
		require.NoError(err)

		err = s.Members().DeleteById(ctx, s.owner.Id)
		require.NoError(err)

		_, err = s.Teams().GetById(ctx, team.Id)
		require.NoError(err)

		_, err = s.Memberships().GetById(ctx, team.Id, s.owner.Id)
		require.ErrorIs(horus.ErrNotExist, err)
	})
}

func (s *MemberStoreTestSuite) TestDeleteByUserIdFromOrg() {
	testCases := []struct {
		desc    string
		org_id  horus.OrgId
		user_id horus.UserId
	}{
		{
			desc:    "org does not exist",
			org_id:  horus.OrgId(uuid.New()),
			user_id: s.owner.UserId,
		},
		{
			desc:    "user does not exist",
			org_id:  s.owner.OrgId,
			user_id: horus.UserId(uuid.New()),
		},
		{
			desc:    "member exists",
			org_id:  s.owner.OrgId,
			user_id: s.owner.UserId,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func(ctx context.Context) {
			require := s.Require()

			err := s.Members().DeleteByUserIdFromOrg(ctx, horus.OrgId(uuid.New()), s.owner.UserId)
			require.NoError(err)
		})
	}
}

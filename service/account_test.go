package service_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus/ent/proto/khepri/horus"
)

type AccountTestSuite struct {
	SuiteWithSilo
}

func TestAccount(t *testing.T) {
	suite.Run(t, &AccountTestSuite{
		SuiteWithSilo: SuiteWithSilo{
			Suite: NewSuiteWithSqliteStore(),
		},
	})
}

func (s *AccountTestSuite) TestGet() {
	s.Run("myself as an owner", func() {
		require := s.Require()

		res, err := s.svc.Account().Get(s.ctx, &horus.GetAccountRequest{
			Id: s.frame.ActingAccount.ID[:],
		})
		require.NoError(err)
		require.Equal(s.frame.ActingAccount.ID[:], res.Id)
	})

	s.Run("myself as a member", func() {
		require := s.Require()

		res, err := s.svc.Account().Get(s.CtxAmigo(), &horus.GetAccountRequest{
			Id: s.amigo.ActingAccount.ID[:],
		})
		require.NoError(err)
		require.Equal(s.amigo.ActingAccount.ID[:], res.Id)
	})

	s.Run("member in same silo as an owner", func() {
		require := s.Require()

		res, err := s.svc.Account().Get(s.ctx, &horus.GetAccountRequest{
			Id: s.amigo.ActingAccount.ID[:],
		})
		require.NoError(err)
		require.Equal(s.amigo.ActingAccount.ID[:], res.Id)
	})

	s.Run("member in same silo as a member", func() {
		require := s.Require()

		res, err := s.svc.Account().Get(s.CtxAmigo(), &horus.GetAccountRequest{
			Id: s.buddy.ActingAccount.ID[:],
		})
		require.NoError(err)
		require.Equal(s.buddy.ActingAccount.ID[:], res.Id)
	})

	s.Run("outsider", func() {
		require := s.Require()

		_, err := s.svc.Account().Get(s.ctx, &horus.GetAccountRequest{
			Id: s.other.ActingAccount.ID[:],
		})
		require.Equal(codes.NotFound, status.Code(err))
	})
}

func (s *AccountTestSuite) TestList() {
	s.Run("list only relevant silos.", func() {
		require := s.Require()

		v1, err := s.svc.Silo().Create(s.ctx, &horus.CreateSiloRequest{
			Silo: &horus.Silo{Alias: "horus1", Name: "Horus"},
		})
		require.NoError(err)

		v2, err := s.svc.Silo().Create(s.ctx, &horus.CreateSiloRequest{
			Silo: &horus.Silo{Alias: "horus2", Name: "Horus"},
		})
		require.NoError(err)

		v3, err := s.svc.Silo().Create(s.CtxAmigo(), &horus.CreateSiloRequest{
			Silo: &horus.Silo{Alias: "khepri", Name: "Khepri"},
		})
		require.NoError(err)

		res, err := s.svc.Account().List(s.ctx, &horus.ListAccountRequest{
			View: horus.ListAccountRequest_BASIC,
		})
		require.NoError(err)

		vs := lo.Map(res.AccountList, func(v *horus.Account, _ int) string { return v.Silo.Alias })
		require.Contains(vs, s.silo.Alias)
		require.Contains(vs, v1.Alias)
		require.Contains(vs, v2.Alias)
		require.NotContains(vs, s.other_silo.Alias)
		require.NotContains(vs, v3.Alias)
	})
}

func (s *AccountTestSuite) TestDelete() {
	s.Run("cannot delete owner", func() {
		require := s.Require()

		_, err := s.svc.Account().Delete(s.ctx, &horus.DeleteAccountRequest{
			Id: s.frame.ActingAccount.ID[:],
		})
		require.Equal(codes.PermissionDenied, status.Code(err))
	})

	s.Run("other member in same silo as an owner", func() {
		require := s.Require()

		_, err := s.svc.Account().Delete(s.ctx, &horus.DeleteAccountRequest{
			Id: s.amigo.ActingAccount.ID[:],
		})
		require.NoError(err)
	})

	s.Run("other member in same silo as a member", func() {
		require := s.Require()

		_, err := s.svc.Account().Delete(s.CtxAmigo(), &horus.DeleteAccountRequest{
			Id: s.buddy.ActingAccount.ID[:],
		})
		require.Equal(codes.PermissionDenied, status.Code(err))
	})

	s.Run("myself as a member", func() {
		require := s.Require()

		_, err := s.svc.Account().Delete(s.CtxAmigo(), &horus.DeleteAccountRequest{
			Id: s.amigo.ActingAccount.ID[:],
		})
		require.NoError(err)
	})

	s.Run("outsider", func() {
		require := s.Require()

		_, err := s.svc.Account().Delete(s.CtxOther(), &horus.DeleteAccountRequest{
			Id: s.amigo.ActingAccount.ID[:],
		})
		require.Equal(codes.NotFound, status.Code(err))
	})
}

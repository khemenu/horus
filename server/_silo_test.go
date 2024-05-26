package server_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/silo"
)

type SiloTestSuite struct {
	SuiteWithSilo
}

func TestSilo(t *testing.T) {
	suite.Run(t, &SiloTestSuite{
		SuiteWithSilo: SuiteWithSilo{
			Suite: NewSuiteWithSqliteStore(),
		},
	})
}

func (s *SiloTestSuite) TestCreate() {
	s.Run("also create an owner account", func() {
		require := s.Require()

		v, err := s.svc.Silo().Create(s.ctx, &horus.CreateSiloRequest{
			Silo: &horus.Silo{
				Alias: "horus",
				Name:  "Horus",
			},
		})
		require.NoError(err)
		require.Equal("horus", v.Alias)
		require.Equal("Horus", v.Name)

		v, err = s.svc.Silo().Get(s.ctx, &horus.GetSiloRequest{
			Id: v.Id,
		})
		require.NoError(err)
		require.Equal("horus", v.Alias)
		require.Equal("Horus", v.Name)

		founder, err := s.frame.Actor.QueryAccounts().
			Where(account.HasSiloWith(silo.ID(uuid.UUID(v.Id)))).
			Only(s.ctx)
		require.NoError(err)
		require.Equal(account.RoleOWNER, founder.Role)
	})

	s.Run("cannot create duplicated alias", func() {
		require := s.Require()

		_, err := s.svc.Silo().Create(s.ctx, &horus.CreateSiloRequest{
			Silo: &horus.Silo{
				Alias: s.silo.Alias,
				Name:  "Horus",
			},
		})
		require.Equal(codes.AlreadyExists, status.Code(err))
	})
}

func (s *SiloTestSuite) TestGet() {
	s.Run("as an owner", func() {
		require := s.Require()

		res, err := s.svc.Silo().Get(s.ctx, &horus.GetSiloRequest{
			Id: s.silo.ID[:],
		})
		require.NoError(err)
		require.Equal(s.silo.Alias, res.Alias)
	})

	s.Run("as a member", func() {
		require := s.Require()

		res, err := s.svc.Silo().Get(s.CtxAmigo(), &horus.GetSiloRequest{
			Id: s.silo.ID[:],
		})
		require.NoError(err)
		require.Equal(s.silo.Alias, res.Alias)
	})

	s.Run("as an outsider", func() {
		require := s.Require()

		_, err := s.svc.Silo().Get(s.ctx, &horus.GetSiloRequest{
			Id: s.other_silo.ID[:],
		})
		require.Equal(codes.NotFound, status.Code(err))
	})
}

func (s *SiloTestSuite) TestUpdate() {
	s.Run("as an owner", func() {
		require := s.Require()

		v, err := s.svc.Silo().Get(s.ctx, &horus.GetSiloRequest{Id: s.silo.ID[:]})
		require.NoError(err)
		require.NotEqual("Khepri", v.Name)

		v.Name = "Khepri"
		_, err = s.svc.Silo().Update(s.ctx, &horus.UpdateSiloRequest{Silo: v})
		require.NoError(err)

		v, err = s.svc.Silo().Get(s.ctx, &horus.GetSiloRequest{Id: v.Id})
		require.NoError(err)
		require.Equal("Khepri", v.Name)
	})

	s.Run("as a member", func() {
		require := s.Require()

		v, err := s.svc.Silo().Get(s.ctx, &horus.GetSiloRequest{Id: s.silo.ID[:]})
		require.NoError(err)

		_, err = s.svc.Silo().Update(s.CtxAmigo(), &horus.UpdateSiloRequest{Silo: v})
		require.Equal(codes.PermissionDenied, status.Code(err))
	})

	s.Run("as an outsider", func() {
		require := s.Require()

		v, err := s.svc.Silo().Get(s.ctx, &horus.GetSiloRequest{Id: s.silo.ID[:]})
		require.NoError(err)

		_, err = s.svc.Silo().Update(s.CtxOther(), &horus.UpdateSiloRequest{Silo: v})
		require.Equal(codes.NotFound, status.Code(err))
	})
}

package service_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus/ent/proto/khepri/horus"
)

type TeamTestSuite struct {
	SuiteWithTeam
}

func TestTeam(t *testing.T) {
	suite.Run(t, &TeamTestSuite{
		SuiteWithTeam: SuiteWithTeam{
			SuiteWithSilo: SuiteWithSilo{
				Suite: NewSuiteWithSqliteStore(),
			},
		},
	})
}

func (s *TeamTestSuite) TestCreate() {
	s.Run("as silo owner", func() {
		require := s.Require()

		v, err := s.svc.Team().Create(s.ctx, &horus.CreateTeamRequest{
			Team: &horus.Team{
				Alias: "crazy-88",
				Name:  "Crazy 88",
				Silo:  &horus.Silo{Id: s.silo.ID[:]},
			},
		})
		require.NoError(err)
		require.Equal("crazy-88", v.Alias)
		require.Equal("Crazy 88", v.Name)
	})

	s.Run("as silo member", func() {
		require := s.Require()

		_, err := s.svc.Team().Create(s.CtxAmigo(), &horus.CreateTeamRequest{
			Team: &horus.Team{
				Alias: "crazy-88",
				Name:  "Crazy 88",
				Silo:  &horus.Silo{Id: s.silo.ID[:]},
			},
		})
		require.Equal(codes.PermissionDenied, status.Code(err))
	})

	s.Run("as outsider", func() {
		require := s.Require()

		_, err := s.svc.Team().Create(s.CtxOther(), &horus.CreateTeamRequest{
			Team: &horus.Team{
				Alias: "crazy-88",
				Name:  "Crazy 88",
				Silo:  &horus.Silo{Id: s.silo.ID[:]},
			},
		})
		require.Equal(codes.NotFound, status.Code(err), err)
	})

	s.Run("with existing alias", func() {
		require := s.Require()

		_, err := s.svc.Team().Create(s.ctx, &horus.CreateTeamRequest{
			Team: &horus.Team{
				Alias: s.team.Alias,
				Name:  "Crazy 42",
				Silo:  &horus.Silo{Id: s.silo.ID[:]},
			},
		})
		require.Equal(codes.AlreadyExists, status.Code(err))
	})

	s.Run("with existing name", func() {
		require := s.Require()

		_, err := s.svc.Team().Create(s.ctx, &horus.CreateTeamRequest{
			Team: &horus.Team{
				Alias: "crazy-88",
				Name:  s.team.Name,
				Silo:  &horus.Silo{Id: s.silo.ID[:]},
			},
		})
		require.NoError(err)
	})
}

func (s *TeamTestSuite) TestGet() {
	s.Run("as silo owner", func() {
		require := s.Require()

		_, err := s.svc.Team().Get(s.ctx, &horus.GetTeamRequest{
			Id: s.team.ID[:],
		})
		require.NoError(err)
	})

	s.Run("as silo member", func() {
		require := s.Require()

		_, err := s.svc.Team().Get(s.CtxAmigo(), &horus.GetTeamRequest{
			Id: s.team.ID[:],
		})
		require.Equal(codes.PermissionDenied, status.Code(err))
	})
}

package service_test

import (
	"context"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/enttest"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/proto/khepri/horus"
	"khepri.dev/horus/ent/team"
	"khepri.dev/horus/service"
	"khepri.dev/horus/service/frame"
)

func NewSuiteWithSqliteStore() Suite {
	return Suite{
		driver_name: "sqlite3",
		source_name: "file:ent?mode=memory&cache=shared&_fk=1",
	}
}

type Suite struct {
	suite.Suite

	driver_name string
	source_name string

	client *ent.Client
	svc    service.Service

	frame *frame.Frame // Frame of actor.
	other *frame.Frame // User who does not have any relation with the actor.

	ctx context.Context
}

func (s *Suite) CtxOther() context.Context {
	return frame.WithContext(s.ctx, s.other)
}

func (s *Suite) SetupSubTest() {
	c := enttest.Open(
		s.T(), s.driver_name, s.source_name,
		enttest.WithOptions(ent.Log(s.T().Log)),
	)

	s.client = c
	s.svc = service.NewService(c)

	s.frame = frame.New()
	s.other = frame.New()
	s.ctx = frame.WithContext(context.Background(), s.frame)

	var err error
	s.frame.Actor, err = c.User.Create().Save(s.ctx)
	if err != nil {
		panic(err)
	}
	s.other.Actor, err = c.User.Create().Save(s.ctx)
	if err != nil {
		panic(err)
	}
}

func (s *Suite) TearDownSubTest() {
	s.ctx = nil
	s.frame = nil
	s.other = nil
	s.svc = nil

	s.client.Close()
	s.client = nil
}

type SuiteWithSilo struct {
	Suite

	amigo *frame.Frame // User who is in the same silo with the actor.
	buddy *frame.Frame // User who is in the same silo with the actor.

	silo       *ent.Silo // Silo owned by actor.
	other_silo *ent.Silo // Silo owned by actor.
}

func (s *SuiteWithSilo) CtxAmigo() context.Context {
	return frame.WithContext(s.ctx, s.amigo)
}

func (s *SuiteWithSilo) CtxBuddy() context.Context {
	return frame.WithContext(s.ctx, s.buddy)
}

func (s *SuiteWithSilo) SetupSubTest() {
	s.Suite.SetupSubTest()

	// Actor's silo.
	{
		v, err := s.svc.Silo().Create(s.ctx, &horus.CreateSiloRequest{
			Silo: &horus.Silo{
				Alias: "x",
				Name:  "Horus",
			},
		})
		if err != nil {
			panic(err)
		}

		s.frame.ActingAccount, err = s.frame.Actor.QueryAccounts().
			Where(account.SiloID(uuid.UUID(v.Id))).
			Only(s.ctx)
		if err != nil {
			panic(err)
		}

		s.silo, err = s.frame.ActingAccount.QuerySilo().
			Only(s.ctx)
		if err != nil {
			panic(err)
		}
	}

	// Other's silo.
	{
		v, err := s.svc.Silo().Create(s.CtxOther(), &horus.CreateSiloRequest{
			Silo: &horus.Silo{
				Alias: "y",
				Name:  "Isis",
			},
		})
		if err != nil {
			panic(err)
		}

		s.other.ActingAccount, err = s.other.Actor.QueryAccounts().
			Where(account.SiloID(uuid.UUID(v.Id))).
			Only(s.ctx)
		if err != nil {
			panic(err)
		}

		s.other_silo, err = s.other.ActingAccount.QuerySilo().
			Only(s.ctx)
		if err != nil {
			panic(err)
		}
	}

	var err error

	s.amigo = frame.New()
	s.amigo.Actor, err = s.client.User.Create().Save(s.ctx)
	if err != nil {
		panic(err)
	}

	s.amigo.ActingAccount, err = s.client.Account.Create().
		SetAlias("amigo").
		SetName("O-Ren Ishii").
		SetOwner(s.amigo.Actor).
		SetSiloID(s.silo.ID).
		SetRole(account.RoleMEMBER).
		Save(s.ctx)
	if err != nil {
		panic(err)
	}

	s.buddy = frame.New()
	s.buddy.Actor, err = s.client.User.Create().Save(s.ctx)
	if err != nil {
		panic(err)
	}

	s.buddy.ActingAccount, err = s.client.Account.Create().
		SetAlias("buddy").
		SetName("Budd").
		SetOwner(s.buddy.Actor).
		SetSiloID(s.silo.ID).
		SetRole(account.RoleMEMBER).
		Save(s.ctx)
	if err != nil {
		panic(err)
	}
}

func (s *SuiteWithSilo) TearDownSubTest() {
	s.amigo = nil
	s.silo = nil
	s.Suite.TearDownSubTest()
}

type SuiteWithTeam struct {
	SuiteWithSilo

	team       *ent.Team
	membership *ent.Membership
}

func (s *SuiteWithTeam) SetupSubTest() {
	s.SuiteWithSilo.SetupSubTest()
	v, err := s.svc.Team().Create(s.ctx, &horus.CreateTeamRequest{
		Team: &horus.Team{
			Alias: "x",
			Name:  "A-Team",
			Silo:  &horus.Silo{Id: s.silo.ID[:]},
		},
	})
	if err != nil {
		panic(err)
	}

	s.membership, err = s.frame.ActingAccount.QueryMemberships().
		Where(membership.HasTeamWith(team.ID(uuid.UUID(v.Id)))).
		Only(s.ctx)
	if err != nil {
		panic(err)
	}

	s.team, err = s.membership.QueryTeam().
		Only(s.ctx)
	if err != nil {
		panic(err)
	}
}

func (s *SuiteWithTeam) TearDownSubTest() {
	s.team = nil
	s.SuiteWithSilo.TearDownSubTest()
}

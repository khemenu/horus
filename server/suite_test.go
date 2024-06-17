package server_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/enttest"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/team"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/role"
	service "khepri.dev/horus/server"
	"khepri.dev/horus/server/frame"
)

func NewSuiteWithSqliteStore() Suite {
	return Suite{
		driver_name: "sqlite3",
		source_name: "file:ent?mode=memory&cache=shared&_fk=1",
	}
}

type Suite struct {
	suite.Suite
	*require.Assertions

	driver_name string
	source_name string

	client *ent.Client
	svc    horus.Server

	me    *frame.Frame // Frame of actor.
	other *frame.Frame // User who does not have any relation with the actor.

	ctx context.Context
}

func (s *Suite) Run(name string, subtest func()) bool {
	return s.Suite.Run(name, func() {
		s.Assertions = s.Require()
		subtest()
	})
}

func (s *Suite) ErrCode(err error, code codes.Code) {
	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(code, st.Code())
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
	s.svc = service.NewServer(c)

	s.me = frame.New()
	s.other = frame.New()
	s.ctx = frame.WithContext(context.Background(), s.me)

	var err error
	s.me.Actor, err = c.User.Create().Save(s.ctx)
	if err != nil {
		panic(err)
	}
	s.me.Token, err = c.Token.Create().
		SetValue("foo").
		SetType(horus.TokenTypeAccess).
		SetDateExpired(time.Now().Add(time.Hour)).
		SetOwner(s.me.Actor).
		Save(s.ctx)
	if err != nil {
		panic(err)
	}

	s.other.Actor, err = c.User.Create().Save(s.ctx)
	if err != nil {
		panic(err)
	}
	s.other.Token, err = c.Token.Create().
		SetValue("bar").
		SetType(horus.TokenTypeAccess).
		SetDateExpired(time.Now().Add(time.Hour)).
		SetOwner(s.other.Actor).
		Save(s.ctx)
	if err != nil {
		panic(err)
	}
}

func (s *Suite) TearDownSubTest() {
	s.ctx = nil
	s.me = nil
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
			Alias: fx.Addr("x"),
			Name:  fx.Addr("Horus"),
		})
		if err != nil {
			panic(err)
		}

		s.me.ActingAccount, err = s.me.Actor.QueryAccounts().
			Where(account.SiloID(uuid.UUID(v.Id))).
			Only(s.ctx)
		if err != nil {
			panic(err)
		}

		s.silo, err = s.me.ActingAccount.QuerySilo().
			Only(s.ctx)
		if err != nil {
			panic(err)
		}
	}

	// Other's silo.
	{
		v, err := s.svc.Silo().Create(s.CtxOther(), &horus.CreateSiloRequest{
			Alias: fx.Addr("y"),
			Name:  fx.Addr("Isis"),
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
		SetRole(role.Member).
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
		SetRole(role.Member).
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
		Alias: fx.Addr("x"),
		Name:  fx.Addr("A-Team"),
		Silo: &horus.GetSiloRequest{Key: &horus.GetSiloRequest_Id{
			Id: s.silo.ID[:],
		}},
	})
	if err != nil {
		panic(err)
	}

	s.membership, err = s.me.ActingAccount.QueryMemberships().
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

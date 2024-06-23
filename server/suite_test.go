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
	"khepri.dev/horus/server"
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

	db  *ent.Client
	svc horus.Server

	me    *frame.Frame // Frame of actor.
	child *frame.Frame // Frame of actor's child.
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

func (s *Suite) CtxMe() context.Context {
	return frame.WithContext(s.ctx, s.me)

}

func (s *Suite) CtxChild() context.Context {
	return frame.WithContext(s.ctx, s.child)
}

func (s *Suite) CtxOther() context.Context {
	return frame.WithContext(s.ctx, s.other)
}

func (s *Suite) initActor() *frame.Frame {
	var err error
	f := frame.New()

	f.Actor, err = s.db.User.Create().Save(s.ctx)
	s.NoError(err)

	f.Token, err = s.db.Token.Create().
		SetValue("token-" + uuid.NewString()).
		SetType(horus.TokenTypeAccess).
		SetDateExpired(time.Now().Add(time.Hour)).
		SetOwner(f.Actor).
		Save(s.ctx)
	s.NoError(err)

	return f
}

func (s *Suite) SetupSubTest() {
	s.Assertions = s.Require()

	c := enttest.Open(
		s.T(), s.driver_name, s.source_name,
		enttest.WithOptions(ent.Log(s.T().Log)),
	)

	s.db = c
	s.svc = server.NewServer(c)
	s.ctx = context.Background()

	s.me = s.initActor()
	s.child = s.initActor()
	s.other = s.initActor()

	_, err := s.child.Actor.Update().
		SetParentID(s.me.Actor.ID).
		Save(s.ctx)
	s.NoError(err)
}

func (s *Suite) TearDownSubTest() {
	s.me = nil
	s.child = nil
	s.other = nil

	s.svc = nil
	s.ctx = nil

	s.db.Close()
	s.db = nil
}

func (s *Suite) AsUser(id []byte) context.Context {
	u, err := s.db.User.Get(s.ctx, uuid.UUID(id))
	s.NoError(err)

	return frame.WithContext(s.ctx, &frame.Frame{Actor: u})
}

type SuiteWithSilo struct {
	Suite

	silo_owner  *frame.Frame // Actor.
	silo_admin  *frame.Frame // User who is in the same silo with the actor.
	silo_member *frame.Frame // User who is in the same silo with the actor.

	silo       *ent.Silo // Silo owned by the actor.
	other_silo *ent.Silo // Silo owned by the other.
}

func (s *SuiteWithSilo) CtxSiloOwner() context.Context {
	return frame.WithContext(s.ctx, s.silo_owner)
}

func (s *SuiteWithSilo) CtxSiloAdmin() context.Context {
	return frame.WithContext(s.ctx, s.silo_admin)
}

func (s *SuiteWithSilo) CtxSiloMember() context.Context {
	return frame.WithContext(s.ctx, s.silo_member)
}

func (s *SuiteWithSilo) SetupSubTest() {
	s.Suite.SetupSubTest()

	// Actor's silo.
	{
		v, err := s.svc.Silo().Create(s.CtxMe(), &horus.CreateSiloRequest{
			Alias: fx.Addr("x"),
			Name:  fx.Addr("Horus"),
		})
		s.NoError(err)

		s.me.ActingAccount, err = s.me.Actor.QueryAccounts().
			Where(account.SiloID(uuid.UUID(v.Id))).
			Only(s.ctx)
		s.NoError(err)

		s.silo, err = s.me.ActingAccount.QuerySilo().
			Only(s.ctx)
		s.NoError(err)
	}

	// Other's silo.
	{
		v, err := s.svc.Silo().Create(s.CtxOther(), &horus.CreateSiloRequest{
			Alias: fx.Addr("y"),
			Name:  fx.Addr("Isis"),
		})
		s.NoError(err)

		s.other.ActingAccount, err = s.other.Actor.QueryAccounts().
			Where(account.SiloID(uuid.UUID(v.Id))).
			Only(s.ctx)
		s.NoError(err)

		s.other_silo, err = s.other.ActingAccount.QuerySilo().
			Only(s.ctx)
		s.NoError(err)
	}

	var err error

	s.silo_owner = s.me

	s.silo_admin = s.initActor()
	s.silo_admin.ActingAccount, err = s.db.Account.Create().SetSiloID(s.silo.ID).SetOwner(s.silo_admin.Actor).
		SetAlias("admin").SetName("O-Ren Ishii").SetRole(role.Admin).Save(s.ctx)
	s.NoError(err)

	s.silo_member = s.initActor()
	s.silo_member.ActingAccount, err = s.db.Account.Create().SetSiloID(s.silo.ID).SetOwner(s.silo_member.Actor).
		SetAlias("member").SetName("Budd").SetRole(role.Member).Save(s.ctx)
	s.NoError(err)
}

func (s *SuiteWithSilo) TearDownSubTest() {
	s.silo_member = nil
	s.silo_admin = nil
	s.silo_owner = nil
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

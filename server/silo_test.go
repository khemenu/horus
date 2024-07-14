package server_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/frame"
)

type SiloTestSuite struct {
	SuiteWithSilo
}

func TestSilo(t *testing.T) {
	s := SiloTestSuite{
		SuiteWithSilo: SuiteWithSilo{
			Suite: NewSuiteWithSqliteStore(),
		},
	}
	suite.Run(t, &s)
}

type SiloAct struct {
	SiloRole role.Role

	OtherSilo bool

	Fail Fail
	Code codes.Code
}

type SiloTestCtx struct {
	Actor    *frame.Frame
	CtxActor context.Context

	TargetSilo *ent.Silo
}

func (a *SiloAct) Prepare(t *SiloTestSuite) *SiloTestCtx {
	actor := t.silo_admin
	err := t.SetSiloRole(t.silo_owner, actor, a.SiloRole)
	t.NoError(err)

	target_silo := t.silo
	if a.OtherSilo {
		target_silo = t.other_silo
	}

	return &SiloTestCtx{
		Actor:      actor,
		CtxActor:   frame.WithContext(t.ctx, actor),
		TargetSilo: target_silo,
	}
}

func (t *SiloTestSuite) baseActs() []SiloAct {
	v := []SiloAct{}
	for _, silo_role := range role.Values() {
		v = append(v, SiloAct{
			SiloRole: silo_role,
		})
	}

	return v
}

func (t *SiloTestSuite) modificationActs() []SiloAct {
	v := []SiloAct{}
	v = append(v, fx.MapV(t.baseActs(), func(act SiloAct) SiloAct {
		act.Fail = act.SiloRole != role.Owner
		act.Code = codes.PermissionDenied
		return act
	})...)
	v = append(v, fx.MapV(t.baseActs(), func(act SiloAct) SiloAct {
		act.OtherSilo = true

		act.Fail = true
		act.Code = codes.NotFound
		return act
	})...)

	return v
}

func (t *SiloTestSuite) TestCreate() {
	t.Run("when a silo is created, an account with owner role is also created", func() {
		v, err := t.svc.Silo().Create(t.CtxSiloOwner(), &horus.CreateSiloRequest{
			Alias:       fx.Addr("monsters"),
			Name:        fx.Addr("Monsters, Inc."),
			Description: fx.Addr("Boo"),
		})
		t.NoError(err)
		t.Equal("monsters", v.Alias)
		t.Equal("Monsters, Inc.", v.Name)
		t.Equal("Boo", v.Description)

		v, err = t.svc.Silo().Get(t.CtxSiloOwner(), horus.SiloByIdV(v.Id))
		t.NoError(err)
		t.Equal("monsters", v.Alias)
		t.Equal("Monsters, Inc.", v.Name)
		t.Equal("Boo", v.Description)

		w, err := t.svc.Account().Get(t.CtxSiloOwner(), horus.AccountByAliasInSilo("founder", horus.SiloByIdV(v.Id)))
		t.NoError(err)
		t.Equal(horus.Role_ROLE_OWNER, w.Role)
	})
	t.Run("silo alias cannot be duplicated", func() {
		_, err := t.svc.Silo().Create(t.CtxMe(), &horus.CreateSiloRequest{
			Alias: &t.silo.Alias,
		})
		t.ErrCode(err, codes.AlreadyExists)
	})
}

func (t *SiloTestSuite) TestGet() {
	acts := []SiloAct{}
	acts = append(acts, t.baseActs()...)
	acts = append(acts, fx.MapV(t.baseActs(), func(act SiloAct) SiloAct {
		act.OtherSilo = true

		act.Fail = true
		act.Code = codes.NotFound
		return act
	})...)

	for _, act := range acts {
		title := fmt.Sprintf("silo %s be retrieved by ", act.Fail)
		title += fx.Cond(
			act.OtherSilo,
			fmt.Sprintf("%s of another silo", act.SiloRole),
			fmt.Sprintf("its %s", act.SiloRole),
		)

		t.Run(title, func() {
			c := act.Prepare(t)

			v, err := t.svc.Silo().Get(c.CtxActor, horus.SiloById(c.TargetSilo.ID))
			if act.Fail {
				t.ErrCode(err, act.Code)
			} else {
				t.NoError(err)
				t.Equal(c.TargetSilo.Alias, v.Alias)
			}
		})
	}

	t.Run("silo that does not exist cannot be retrieved", func() {
		_, err := t.svc.Silo().Get(t.CtxOther(), horus.SiloByAlias("not exist"))
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *SiloTestSuite) TestUpdate() {
	for _, act := range t.modificationActs() {
		title := fmt.Sprintf("silo %s be updated by ", act.Fail)
		title += fx.Cond(
			act.OtherSilo,
			fmt.Sprintf("%s of another silo", act.SiloRole),
			fmt.Sprintf("its %s", act.SiloRole),
		)

		t.Run(title, func() {
			c := act.Prepare(t)

			v, err := t.svc.Silo().Update(c.CtxActor, &horus.UpdateSiloRequest{
				Key:         horus.SiloById(c.TargetSilo.ID),
				Alias:       fx.Addr("monsters"),
				Name:        fx.Addr("Monsters, Inc."),
				Description: fx.Addr("Boo"),
			})
			if act.Fail {
				t.ErrCode(err, act.Code)
			} else {
				t.NoError(err)
				t.Equal("monsters", v.Alias)
				t.Equal("Monsters, Inc.", v.Name)
				t.Equal("Boo", v.Description)

				v, err = t.svc.Silo().Get(c.CtxActor, horus.SiloById(c.TargetSilo.ID))
				t.NoError(err)
				t.Equal("monsters", v.Alias)
				t.Equal("Monsters, Inc.", v.Name)
				t.Equal("Boo", v.Description)
			}
		})
	}

	t.Run("silo cannot be updated if the silo does not exist", func() {
		_, err := t.svc.Silo().Update(t.CtxSiloOwner(), &horus.UpdateSiloRequest{
			Key:   horus.SiloByAlias("not exist"),
			Alias: fx.Addr("monsters"),
		})
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *SiloTestSuite) TestDelete() {
	for _, act := range t.modificationActs() {
		title := fmt.Sprintf("silo %s be deleted by ", act.Fail)
		title += fx.Cond(
			act.OtherSilo,
			fmt.Sprintf("%s of another silo", act.SiloRole),
			fmt.Sprintf("its %s", act.SiloRole),
		)

		t.Run(title, func() {
			c := act.Prepare(t)
			_, err := t.svc.Silo().Delete(c.CtxActor, horus.SiloById(c.TargetSilo.ID))
			if act.Fail {
				t.ErrCode(err, act.Code)
			} else {
				t.NoError(err)

				_, err = t.svc.Silo().Get(c.CtxActor, horus.SiloById(c.TargetSilo.ID))
				t.ErrCode(err, codes.NotFound)
			}
		})
	}

	t.Run("silo cannot be deleted if the silo does not exist", func() {
		_, err := t.svc.Silo().Delete(t.CtxSiloOwner(), horus.SiloByAlias("not exist"))
		t.ErrCode(err, codes.NotFound)
	})
}

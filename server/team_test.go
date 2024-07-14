package server_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/frame"
)

type TeamTestSuite struct {
	SuiteWithTeam
}

func TestTeam(t *testing.T) {
	s := TeamTestSuite{
		SuiteWithTeam: SuiteWithTeam{
			SuiteWithSilo: SuiteWithSilo{
				Suite: NewSuiteWithSqliteStore(),
			},
		},
	}
	suite.Run(t, &s)
}

type TeamAct struct {
	SiloRole role.Role
	TeamRole role.Role

	OtherSilo bool
	OtherTeam bool

	Fail Fail
	Code codes.Code
}

func (a *TeamAct) IsWithTeam() bool {
	return !a.TeamRole.IsNil()
}

type TeamTestCtx struct {
	Actor    *frame.Frame
	CtxActor context.Context

	TargetSilo *ent.Silo
	TargetTeam *ent.Team
}

func (a *TeamAct) Prepare(t *TeamTestSuite) *TeamTestCtx {
	actor := fx.Cond(a.IsWithTeam(), t.team_admin, t.silo_admin)
	err := t.SetSiloRole(t.silo_owner, actor, a.SiloRole)
	t.NoError(err)
	if a.IsWithTeam() {
		err := t.SetTeamRole(t.silo_owner, actor, horus.TeamById(t.team.ID), a.TeamRole)
		t.NoError(err)
	}

	target_silo := t.silo
	if a.OtherSilo {
		target_silo = t.other_silo
	}

	target_team := t.team
	if a.OtherSilo {
		target_team = t.other_team
	} else if a.OtherTeam {
		v, err := t.svc.Team().Create(t.CtxSiloOwner(), &horus.CreateTeamRequest{
			Silo: horus.SiloById(t.silo.ID),
		})
		t.NoError(err)

		target_team, err = t.db.Team.Get(t.ctx, uuid.UUID(v.Id))
		t.NoError(err)
	}

	return &TeamTestCtx{
		Actor:    actor,
		CtxActor: frame.WithContext(t.ctx, actor),

		TargetSilo: target_silo,
		TargetTeam: target_team,
	}
}

func (t *TeamTestSuite) baseActs() []TeamAct {
	v := []TeamAct{}
	for _, silo_role := range role.Values() {
		v = append(v, TeamAct{
			SiloRole: silo_role,
		})
	}
	for _, silo_role := range role.Values() {
		for _, team_role := range role.Values() {
			v = append(v, TeamAct{
				SiloRole: silo_role,
				TeamRole: team_role,
			})
		}
	}

	return v
}

func (t *TeamTestSuite) modificationActs() []TeamAct {
	v := []TeamAct{}
	v = append(v, fx.MapV(t.baseActs(), func(act TeamAct) TeamAct {
		silo_roles := []role.Role{role.Owner, role.Admin}
		team_roles := []role.Role{role.Owner}

		act.Fail = Fail(!fx.Or(
			slices.Contains(silo_roles, act.SiloRole),
			slices.Contains(team_roles, act.TeamRole),
		))
		act.Code = fx.Cond(
			act.TeamRole.IsNil(),
			codes.NotFound, // Silo member who does not have a team.
			codes.PermissionDenied,
		)
		return act
	})...)
	v = append(v, fx.MapV(t.baseActs(), func(act TeamAct) TeamAct {
		act.OtherSilo = true

		act.Fail = true
		act.Code = codes.NotFound
		return act
	})...)
	v = append(v, fx.MapV(t.baseActs()[3:], func(act TeamAct) TeamAct {
		act.OtherTeam = true

		silo_roles := []role.Role{role.Owner, role.Admin}

		act.Fail = Fail(!slices.Contains(silo_roles, act.SiloRole))
		act.Code = codes.NotFound
		return act
	})...)

	return v
}

func (t *TeamTestSuite) TestCreate() {
	acts := []TeamAct{}
	acts = append(acts, fx.MapV(t.baseActs(), func(act TeamAct) TeamAct {
		silo_roles := []role.Role{role.Owner, role.Admin}

		act.Fail = Fail(!slices.Contains(silo_roles, act.SiloRole))
		act.Code = codes.PermissionDenied
		return act
	})...)
	acts = append(acts, fx.MapV(t.baseActs(), func(act TeamAct) TeamAct {
		act.OtherSilo = true

		act.Fail = true
		act.Code = codes.NotFound
		return act
	})...)

	for _, act := range acts {
		title := fmt.Sprintf("team %s be created by ", act.Fail)
		title += fx.Cond(
			act.OtherSilo,
			fmt.Sprintf("%s of another silo", act.SiloRole),
			fmt.Sprintf("the silo %s", act.SiloRole),
		)
		if act.IsWithTeam() {
			title += fmt.Sprintf(" who is also a team %s", act.TeamRole)
		}

		t.Run(title, func() {
			c := act.Prepare(t)

			_, err := t.svc.Team().Create(c.CtxActor, &horus.CreateTeamRequest{
				Silo: horus.SiloById(c.TargetSilo.ID),
			})
			if act.Fail {
				t.ErrCode(err, act.Code)
			} else {
				t.NoError(err)
			}
		})
	}

	t.Run("when a team is created, a membership with owner role is also created", func() {
		s, err := t.svc.Team().Create(t.CtxSiloOwner(), &horus.CreateTeamRequest{
			Silo: horus.SiloById(t.silo.ID),
		})
		t.NoError(err)

		v, err := t.svc.Membership().Get(t.CtxSiloOwner(), horus.MembershipByAccountInTeam(
			horus.AccountById(t.silo_owner.ActingAccount.ID),
			horus.TeamByIdV(s.Id),
		))
		t.NoError(err)
		t.Equal(horus.Role_ROLE_OWNER, v.Role)
	})
	t.Run("team cannot be created if the silo does not exist", func() {
		_, err := t.svc.Team().Create(t.CtxSiloOwner(), &horus.CreateTeamRequest{
			Silo: horus.SiloByAlias("not exist"),
		})
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *TeamTestSuite) TestGet() {
	acts := []TeamAct{}
	acts = append(acts, fx.MapV(t.baseActs(), func(act TeamAct) TeamAct {
		silo_roles := []role.Role{role.Owner, role.Admin}
		team_roles := []role.Role{role.Owner, role.Admin, role.Member}

		act.Fail = Fail(!fx.Or(
			slices.Contains(silo_roles, act.SiloRole),
			slices.Contains(team_roles, act.TeamRole),
		))
		act.Code = codes.NotFound
		return act
	})...)
	acts = append(acts, fx.MapV(t.baseActs(), func(act TeamAct) TeamAct {
		act.OtherSilo = true

		act.Fail = true
		act.Code = codes.NotFound
		return act
	})...)
	acts = append(acts, fx.MapV(t.baseActs()[3:], func(act TeamAct) TeamAct {
		act.OtherTeam = true

		silo_roles := []role.Role{role.Owner, role.Admin}

		act.Fail = Fail(!slices.Contains(silo_roles, act.SiloRole))
		act.Code = codes.NotFound
		return act
	})...)

	for _, act := range acts {
		title := fmt.Sprintf("team %s be retrieved by the ", act.Fail)
		title += fx.Cond(
			act.OtherSilo,
			fmt.Sprintf("%s of another silo", act.SiloRole),
			fmt.Sprintf("silo %s", act.SiloRole),
		)
		if act.IsWithTeam() {
			title += " who is also "
			title += fx.Cond(
				act.OtherTeam,
				fmt.Sprintf("%s of another team in the silo", act.TeamRole),
				fmt.Sprintf("a team %s", act.TeamRole),
			)
		}

		t.Run(title, func() {
			c := act.Prepare(t)

			_, err := t.svc.Team().Get(c.CtxActor, horus.TeamById(c.TargetTeam.ID))
			if act.Fail {
				t.ErrCode(err, act.Code)
			} else {
				t.NoError(err)
			}
		})
	}

	t.Run("team that does not exist cannot be retrieved", func() {
		_, err := t.svc.Team().Get(t.CtxSiloOwner(), horus.TeamByAliasInSilo("not exist", horus.SiloById(t.silo.ID)))
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *TeamTestSuite) TestUpdate() {
	for _, act := range t.modificationActs() {
		title := fmt.Sprintf("team %s be updated by the ", act.Fail)
		title += fx.Cond(
			act.OtherSilo,
			fmt.Sprintf("%s of another silo", act.SiloRole),
			fmt.Sprintf("silo %s", act.SiloRole),
		)
		if act.IsWithTeam() {
			title += " who is also "
			title += fx.Cond(
				act.OtherTeam,
				fmt.Sprintf("%s of another team in the silo", act.TeamRole),
				fmt.Sprintf("a team %s", act.TeamRole),
			)
		}

		t.Run(title, func() {
			c := act.Prepare(t)

			v, err := t.svc.Team().Update(c.CtxActor, &horus.UpdateTeamRequest{
				Key:         horus.TeamById(c.TargetTeam.ID),
				Alias:       fx.Addr("crazy88"),
				Name:        fx.Addr("Crazy 88"),
				Description: fx.Addr("Yakuza"),
			})
			if act.Fail {
				t.ErrCode(err, act.Code)
			} else {
				t.NoError(err)
				t.Equal("crazy88", v.Alias)
				t.Equal("Crazy 88", v.Name)
				t.Equal("Yakuza", v.Description)

				v, err = t.svc.Team().Get(c.CtxActor, horus.TeamById(c.TargetTeam.ID))
				t.NoError(err)
				t.Equal("crazy88", v.Alias)
				t.Equal("Crazy 88", v.Name)
				t.Equal("Yakuza", v.Description)
			}
		})
	}

	t.Run("team cannot be updated if the team does not exist", func() {
		_, err := t.svc.Team().Update(t.CtxSiloOwner(), &horus.UpdateTeamRequest{
			Key:   horus.TeamByAliasInSilo("not exist", horus.SiloById(t.silo.ID)),
			Alias: fx.Addr("crazy88"),
		})
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *TeamTestSuite) TestDelete() {
	for _, act := range t.modificationActs() {
		title := fmt.Sprintf("team %s be deleted by the ", act.Fail)
		title += fx.Cond(
			act.OtherSilo,
			fmt.Sprintf("%s of another silo", act.SiloRole),
			fmt.Sprintf("silo %s", act.SiloRole),
		)
		if act.IsWithTeam() {
			title += " who is also "
			title += fx.Cond(
				act.OtherTeam,
				fmt.Sprintf("%s of another team in the silo", act.TeamRole),
				fmt.Sprintf("a team %s", act.TeamRole),
			)
		}

		t.Run(title, func() {
			c := act.Prepare(t)

			_, err := t.svc.Team().Delete(c.CtxActor, horus.TeamById(c.TargetTeam.ID))
			if act.Fail {
				t.ErrCode(err, act.Code)
			} else {
				t.NoError(err)

				_, err = t.svc.Team().Get(c.CtxActor, horus.TeamById(c.TargetTeam.ID))
				t.ErrCode(err, codes.NotFound)
			}
		})
	}

	t.Run("team cannot be deleted if the team does not exist", func() {
		_, err := t.svc.Team().Delete(t.CtxSiloOwner(), horus.TeamByAliasInSilo("not exist", horus.SiloById(t.silo.ID)))
		t.ErrCode(err, codes.NotFound)
	})
}

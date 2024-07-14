package server_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
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

func (t *TeamTestSuite) TestCreate() {
	type Act struct {
		SiloRole role.Role
		TeamRole role.Role

		OtherSilo bool

		Fail Fail
		Code codes.Code
	}

	base := []Act{}
	for _, silo_role := range role.Values() {
		base = append(base, Act{
			SiloRole: silo_role,
		})
	}
	for _, silo_role := range role.Values() {
		for _, team_role := range role.Values() {
			base = append(base, Act{
				SiloRole: silo_role,
				TeamRole: team_role,
			})
		}
	}

	acts := []Act{}
	acts = append(acts, fx.MapV(base, func(act Act) Act {
		silo_roles := []role.Role{role.Owner, role.Admin}

		act.Fail = Fail(!slices.Contains(silo_roles, act.SiloRole))
		act.Code = codes.PermissionDenied
		return act
	})...)
	acts = append(acts, fx.MapV(base, func(act Act) Act {
		act.OtherSilo = true

		act.Fail = true
		act.Code = codes.NotFound
		return act
	})...)

	for _, act := range acts {
		has_team := !act.TeamRole.IsNil()

		title := fmt.Sprintf("silo %s ", act.SiloRole)
		if has_team {
			title += fmt.Sprintf("who is also a team %s ", act.TeamRole)
		}
		title += fmt.Sprintf("%s create a team", act.Fail)
		if act.OtherSilo {
			title += " in another silo"
		}

		t.Run(title, func() {
			actor := fx.Cond(has_team, t.team_admin, t.silo_admin)
			ctx := frame.WithContext(t.ctx, actor)

			err := t.SetSiloRole(t.silo_owner, actor, act.SiloRole)
			t.NoError(err)
			if has_team {
				err := t.SetTeamRole(t.silo_owner, actor, horus.TeamById(t.team.ID), act.TeamRole)
				t.NoError(err)
			}

			target_silo := fx.Cond(act.OtherSilo, t.other_silo, t.silo)

			_, err = t.svc.Team().Create(ctx, &horus.CreateTeamRequest{
				Silo: horus.SiloById(target_silo.ID),
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
	type Act struct {
		SiloRole role.Role
		TeamRole role.Role

		OtherSilo bool
		OtherTeam bool

		Fail Fail
		Code codes.Code
	}

	base := []Act{}
	for _, silo_role := range role.Values() {
		base = append(base, Act{
			SiloRole: silo_role,
		})
	}
	for _, silo_role := range role.Values() {
		for _, team_role := range role.Values() {
			base = append(base, Act{
				SiloRole: silo_role,
				TeamRole: team_role,
			})
		}
	}

	base_with_team := base[3:]

	acts := []Act{}
	acts = append(acts, fx.MapV(base, func(act Act) Act {
		silo_roles := []role.Role{role.Owner, role.Admin}
		team_roles := []role.Role{role.Owner, role.Admin, role.Member}

		act.Fail = Fail(!fx.Or(
			slices.Contains(silo_roles, act.SiloRole),
			slices.Contains(team_roles, act.TeamRole),
		))
		act.Code = codes.NotFound
		return act
	})...)
	acts = append(acts, fx.MapV(base, func(act Act) Act {
		act.OtherSilo = true

		act.Fail = true
		act.Code = codes.NotFound
		return act
	})...)
	acts = append(acts, fx.MapV(base_with_team, func(act Act) Act {
		act.OtherTeam = true

		silo_roles := []role.Role{role.Owner, role.Admin}

		act.Fail = Fail(!slices.Contains(silo_roles, act.SiloRole))
		act.Code = codes.NotFound
		return act
	})...)

	for _, act := range acts {
		has_team := !act.TeamRole.IsNil()

		title := fmt.Sprintf("silo %s ", act.SiloRole)
		if has_team {
			title += fmt.Sprintf("who is also a team %s ", act.TeamRole)
		}
		title += fmt.Sprintf("%s retrieve ", act.Fail)
		if act.OtherSilo {
			title += "a team in another silo"
		} else if act.OtherTeam {
			title += "another team in the silo"
		} else if has_team {
			title += "their team"
		} else {
			title += "a team in their silo"
		}

		t.Run(title, func() {
			actor := fx.Cond(has_team, t.team_admin, t.silo_admin)
			ctx := frame.WithContext(t.ctx, actor)

			err := t.SetSiloRole(t.silo_owner, actor, act.SiloRole)
			t.NoError(err)
			if has_team {
				err := t.SetTeamRole(t.silo_owner, actor, horus.TeamById(t.team.ID), act.TeamRole)
				t.NoError(err)
			}

			target_team_id := t.team.ID[:]
			if act.OtherSilo {
				target_team_id = t.other_team.ID[:]
			} else if act.OtherTeam {
				v, err := t.svc.Team().Create(t.CtxSiloOwner(), &horus.CreateTeamRequest{
					Silo: horus.SiloById(t.silo.ID),
				})
				t.NoError(err)
				target_team_id = v.Id
			}

			_, err = t.svc.Team().Get(ctx, horus.TeamByIdV(target_team_id))
			if act.Fail {
				t.ErrCode(err, act.Code)
			} else {
				t.NoError(err)
			}
		})
	}

	t.Run("team cannot be retrieved if the team does not exist", func() {
		_, err := t.svc.Team().Get(t.CtxSiloOwner(), horus.TeamByAliasInSilo("not exist", horus.SiloById(t.silo.ID)))
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *TeamTestSuite) TestModify() {
	type Act struct {
		SiloRole role.Role
		TeamRole role.Role

		OtherSilo bool
		OtherTeam bool

		Fail Fail
		Code codes.Code
	}

	base := []Act{}
	for _, silo_role := range role.Values() {
		base = append(base, Act{
			SiloRole: silo_role,
		})
	}
	for _, silo_role := range role.Values() {
		for _, team_role := range role.Values() {
			base = append(base, Act{
				SiloRole: silo_role,
				TeamRole: team_role,
			})
		}
	}

	base_with_team := base[3:]

	acts := []Act{}
	acts = append(acts, fx.MapV(base, func(act Act) Act {
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
	acts = append(acts, fx.MapV(base, func(act Act) Act {
		act.OtherSilo = true

		act.Fail = true
		act.Code = codes.NotFound
		return act
	})...)
	acts = append(acts, fx.MapV(base_with_team, func(act Act) Act {
		act.OtherTeam = true

		silo_roles := []role.Role{role.Owner, role.Admin}

		act.Fail = Fail(!slices.Contains(silo_roles, act.SiloRole))
		act.Code = codes.NotFound
		return act
	})...)

	for _, act := range acts {
		has_team := !act.TeamRole.IsNil()

		title := fmt.Sprintf("silo %s ", act.SiloRole)
		if has_team {
			title += fmt.Sprintf("who is also a team %s ", act.TeamRole)
		}
		title += fmt.Sprintf("%s modify ", act.Fail)
		if act.OtherSilo {
			title += "a team in another silo"
		} else if act.OtherTeam {
			title += "another team in the silo"
		} else if has_team {
			title += "their team"
		} else {
			title += "a team in their silo"
		}

		t.Run(title, func() {
			actor := fx.Cond(has_team, t.team_admin, t.silo_admin)
			ctx := frame.WithContext(t.ctx, actor)

			err := t.SetSiloRole(t.silo_owner, actor, act.SiloRole)
			t.NoError(err)
			if has_team {
				err := t.SetTeamRole(t.silo_owner, actor, horus.TeamById(t.team.ID), act.TeamRole)
				t.NoError(err)
			}

			target_team_id := t.team.ID[:]
			if act.OtherSilo {
				target_team_id = t.other_team.ID[:]
			} else if act.OtherTeam {
				v, err := t.svc.Team().Create(t.CtxSiloOwner(), &horus.CreateTeamRequest{
					Silo: horus.SiloById(t.silo.ID),
				})
				t.NoError(err)
				target_team_id = v.Id
			}

			_, err = t.svc.Team().Update(ctx, &horus.UpdateTeamRequest{
				Key:   horus.TeamByIdV(target_team_id),
				Alias: fx.Addr("crazy88"),
				Name:  fx.Addr("Crazy 88"),
			})
			if act.Fail {
				t.ErrCode(err, act.Code)
			} else {
				t.NoError(err)
			}

			_, err = t.svc.Team().Delete(ctx, horus.TeamByIdV(target_team_id))
			if act.Fail {
				t.ErrCode(err, act.Code)
			} else {
				t.NoError(err)
			}
		})
	}

	t.Run("team cannot be modified if the team does not exist", func() {
		_, err := t.svc.Team().Update(t.CtxSiloOwner(), &horus.UpdateTeamRequest{
			Key:   horus.TeamByAliasInSilo("not exist", horus.SiloById(t.silo.ID)),
			Alias: fx.Addr("crazy88"),
		})
		t.ErrCode(err, codes.NotFound)

		_, err = t.svc.Team().Delete(t.CtxSiloOwner(), horus.TeamByAliasInSilo("not exist", horus.SiloById(t.silo.ID)))
		t.ErrCode(err, codes.NotFound)
	})
}

package server_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/frame"
)

type TeamTestSuite struct {
	SuiteWithSilo
}

func TestTeam(t *testing.T) {
	s := TeamTestSuite{
		SuiteWithSilo: SuiteWithSilo{
			Suite: NewSuiteWithSqliteStore(),
		},
	}
	suite.Run(t, &s)
}

func (t *TeamTestSuite) TestCreate() {
	for _, act := range []struct {
		Actor role.Role
		Fail  bool
	}{
		{
			Actor: role.Owner,
		},
		{
			Actor: role.Admin,
		},
		{
			Actor: role.Member,
			Fail:  true,
		},
	} {
		title := "silo " + strings.ToLower(string(act.Actor)) + " "
		title += fx.Cond(act.Fail, "cannot", "can")
		title += " create a team"

		t.Run(title, func() {
			actor := t.silo_admin
			err := t.SetSiloRole(t.silo_owner, actor, act.Actor)
			t.NoError(err)

			ctx := frame.WithContext(t.ctx, actor)
			_, err = t.svc.Team().Create(ctx, &horus.CreateTeamRequest{
				Silo: horus.SiloById(t.silo.ID),
			})
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
			} else {
				t.NoError(err)
			}
		})

		for _, team_role := range []role.Role{
			role.Owner,
			role.Admin,
			role.Member,
		} {
			title := "team " + strings.ToLower(string(team_role)) + " "
			title += "who is also a silo " + strings.ToLower(string(act.Actor)) + " "
			title += fx.Cond(act.Fail, "cannot", "can")
			title += " create a team"

			// t.Run(title, func() {
			// 	actor := t.silo_admin
			// 	err := t.setRole(t.silo_owner, actor, act.Actor)
			// 	t.NoError(err)

			// })
		}
	}
}

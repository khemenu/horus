package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
)

var CmdCreateTeam = &cli.Command{
	Name:      "team",
	Args:      true,
	ArgsUsage: " [TEAM_ALIAS] in <SILO_ID>",

	Action: func(ctx *cli.Context) error {
		var (
			team_alias string
			silo_id    string
		)
		if ctx.Args().Len() < 2 {
			return fmt.Errorf(`requires at least 2 arguments, e.g. "in" and <SILO_ID>`)
		}
		if ctx.Args().Len() > 3 {
			return fmt.Errorf(`accepts up to 3 arguments`)
		}

		n := 0
		if ctx.Args().Len() == 3 {
			n++
			team_alias = ctx.Args().Get(0)
		}
		if p := ctx.Args().Get(n); p != "in" {
			return fmt.Errorf(`expected a preposition "in" but found %s`, p)
		}
		silo_id = ctx.Args().Get(n + 1)

		conf := ConfFrom(ctx.Context)
		if err := conf.Client.notToBeBareServe(); err != nil {
			return err
		}

		c, err := conf.Client.connect(ctx.Context)
		if err != nil {
			return err
		}

		silo_uuid, _ := uuid.Parse(silo_id)
		v, err := c.Team().Create(ctx.Context, &horus.CreateTeamRequest{
			Alias: &team_alias,
			Silo: &horus.Silo{
				Id:    silo_uuid[:],
				Alias: silo_id,
			},
		})
		if err != nil {
			return err
		}

		o := fmt.Sprintf("%s %s", uuid.UUID(v.Id), v.Alias)
		return conf.Reporter.Report(v, o)
	},
}

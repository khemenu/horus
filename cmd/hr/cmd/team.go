package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/hr/env"
)

var CmdCreateTeam = &cli.Command{
	Name:      "team",
	Args:      true,
	ArgsUsage: " [TEAM_ALIAS] in <SILO_ID>",

	Action: func(ctx_ *cli.Context) error {
		ctx := ctx_.Context
		args := ctx_.Args()

		off := 0
		switch args.Len() {
		case 2:
			off++
		case 3:

		default:
			return fmt.Errorf("requires exactly 2 or 3 arguments")
		}
		if p := args.Get(off + 1); p != "in" {
			return fmt.Errorf(`expected a preposition "in" but found %s`, p)
		}

		var (
			team_alias = args.Get(off + 0)
			silo_id    = args.Get(off + 2)
		)

		e := env.From(ctx)
		h, err := e.Connect(ctx)
		if err != nil {
			return err
		}

		req := &horus.CreateTeamRequest{
			Silo: horus.SiloByQuery(silo_id),
		}
		if team_alias != "" {
			req.Alias = &team_alias
		}

		v, err := h.Team().Create(ctx, req)
		if err != nil {
			return fmt.Errorf("execute: %w", err)
		}

		o := uuid.UUID(v.Id).String()
		return e.Report(v, o)
	},
}

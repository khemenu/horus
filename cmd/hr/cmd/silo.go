package cmd

import (
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
)

var CmdCreateSilo = &cli.Command{
	Name:      "silo",
	Args:      true,
	ArgsUsage: " [SILO_ALIAS]",
	Action: func(ctx *cli.Context) error {
		var alias string
		if ctx.Args().Len() > 0 {
			alias = ctx.Args().Get(0)
		}

		conf := ConfFrom(ctx.Context)
		c, err := conf.Client.connect(ctx.Context)
		if err != nil {
			return err
		}

		v, err := c.Silo().Create(ctx.Context, &horus.CreateSiloRequest{
			Alias: &alias,
		})
		if err != nil {
			return err
		}

		o := uuid.UUID(v.Id).String()
		return conf.Reporter.Report(v, o)
	},
}

package cmd

import (
	"github.com/urfave/cli/v2"
	"khepri.dev/horus/cmd/conf"
)

var Commands = []*cli.Command{
	CmdServe,
	conf.CmdVersion,
}

var App = &cli.App{
	Name:  "horus",
	Usage: "Horus server",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c"},
			Value:   "horus.yaml",
			Usage:   "path to a config file",
		},
	},
	Before: func(ctx *cli.Context) error {
		_, err := conf.InitCmd(ctx, func(c *conf.Config) {})
		if err != nil {
			return err
		}

		return nil
	},

	Commands: Commands,
}

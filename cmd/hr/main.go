package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus/cmd/hr/cmd"
	"khepri.dev/horus/log"
)

func main() {
	var conf *cmd.Config

	app := &cli.App{
		Name:        "hr",
		Description: "Horus client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Value:   "horus.yaml",
				Usage:   "path to a config file",
			},
		},
		Before: func(ctx *cli.Context) error {
			p := ctx.String("conf")
			var err error
			conf, err = cmd.ReadConfig(p)
			if err != nil {
				return fmt.Errorf("read config: %w", err)
			}

			l := conf.Log.NewLogger()
			if conf.Client.Db.UseBare {
				l.Warn("use bare service")
			}

			ctx.Context = log.Into(ctx.Context, l)
			ctx.Context = cmd.ConfInto(ctx.Context, conf)
			return nil
		},

		Commands: cmd.Commands,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

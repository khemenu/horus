package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus/cmd/conf"
	"khepri.dev/horus/cmd/hr/env"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/log"
)

var App = &cli.App{
	Name:        "hr",
	Description: "Horus client",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c"},
			Value:   "horus.yaml",
			Usage:   "path to a config file",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "formatting output",
		},
		&cli.BoolFlag{
			Name:  "no-log",
			Value: false,
			Usage: "disable logging",
		},
		&cli.StringFlag{
			Name:  "as",
			Usage: "act as a given user",

			Action: func(ctx_ *cli.Context, s string) error {
				if s == "" {
					return nil
				}

				ctx := ctx_.Context
				c := conf.From(ctx)
				l := log.From(ctx)
				e := env.From(ctx)

				e.ActorId = s

				if t := c.Hr.Connect.With; t != "db" {
					c.Hr.Connect.With = "db"
					l.Warn(`connection type is changed to "db" because an actor is being used`, slog.String("was", t))
				}

				return nil
			},
		},
	},
	Before: func(ctx *cli.Context) error {
		_, err := conf.InitCmd(ctx, func(c *conf.Config) {})
		if err != nil {
			return err
		}

		ctx.Context = env.Into(ctx.Context, &env.Env{})
		return nil
	},
	After: func(ctx_ *cli.Context) error {
		ctx := ctx_.Context

		e := env.From(ctx)
		return e.CleanUp(ctx)
	},

	Commands: Commands,

	ExitErrHandler: func(ctx *cli.Context, err error) {
		if err == nil {
			return
		}

		code := 1
		if s, ok := status.FromError(err); ok {
			code = int(s.Code())
		} else if ent.IsNotFound(err) {
			code = int(codes.NotFound)
		}

		log.From(ctx.Context).Error(
			"exit with error",
			slog.Int("code", code),
			slog.String("err", err.Error()),
		)

		fmt.Println(err.Error())
		os.Exit(code)
	},
}

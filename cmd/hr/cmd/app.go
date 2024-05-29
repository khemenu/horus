package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/metadata"
	"khepri.dev/horus/ent/user"
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

			Action: func(ctx *cli.Context, s string) error {
				if s == "" {
					return nil
				}

				conf := ConfFrom(ctx.Context)
				conf.Client.has_actor = true
				if conf.Client.ConnectWith != "db" {
					conf.Client.ConnectWith = "db"
					conf.Log.NewLogger().Warn(`connect with DB as actor is provided`)
				}

				return nil
			},
		},
	},
	Before: func(ctx *cli.Context) error {
		conf_path := ctx.String("conf")
		conf, err := ReadConfig(conf_path)
		if err != nil {
			return fmt.Errorf("read config: %w", err)
		}
		if enabled := ctx.Bool("no-log"); enabled {
			conf.Log.Enabled = &enabled
		}
		if v := ctx.String("output"); v != "" {
			entries := strings.SplitN(v, "=", 2)
			if len(entries) == 0 {
				return fmt.Errorf("output format not specified")
			}

			format := entries[0]
			var opt string
			if len(entries) > 1 {
				opt = entries[1]
			}
			conf.Reporter.Format = format
			switch format {
			case "template":
				conf.Reporter.Template = opt

			default:
				// Invalid format is handled by `NewReporter`
			}
		}
		if err := conf.Evaluate(); err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}

		l := conf.Log.NewLogger()
		ctx.Context = log.Into(ctx.Context, l)
		ctx.Context = ConfInto(ctx.Context, conf)

		{
			var opt string
			switch conf.Reporter.Format {
			case "template":
				opt = conf.Reporter.Template
			}

			conf.Reporter.reporter, err = NewReporter(conf.Reporter.Format, opt)
			if err != nil {
				return fmt.Errorf("invalid reporter config: %w", err)
			}
		}

		if actor_id := ctx.String("as"); actor_id != "" {
			conf.Client.has_actor = true
			if conf.Client.ConnectWith != "db" {
				conf.Client.ConnectWith = "db"
				conf.Log.NewLogger().Warn(`connect with DB as actor is provided`)
			}

			if _, err := conf.Client.connect(ctx.Context); err != nil {
				return fmt.Errorf("connect server: %w", err)
			}

			if _, err := uuid.Parse(actor_id); err != nil {
				user_id, err := conf.Client.db.User.Query().
					Where(user.AliasEQ(actor_id)).
					OnlyID(ctx.Context)
				if err != nil {
					return fmt.Errorf("actor not found: %w", err)
				}

				actor_id = user_id.String()
			}
			ctx.Context = metadata.AppendToOutgoingContext(ctx.Context, "actor-uuid", actor_id)
		}

		if err := conf.Evaluate(); err != nil {
			panic("invalid config modification")
		}
		return nil
	},
	After: func(ctx *cli.Context) error {
		conf := ConfFrom(ctx.Context)
		l := log.From(ctx.Context)
		if err := conf.Client.CleanUp(ctx.Context); err != nil {
			l.Error("failed to clean up the client", slog.String("err", err.Error()))
		}
		return nil
	},

	Commands: Commands,

	ExitErrHandler: func(ctx *cli.Context, err error) {
		if err == nil {
			return
		}

		conf := ConfFrom(ctx.Context)
		l := log.From(ctx.Context)

		exit_err := conf.Reporter.ExitWithErr(err)
		l.Error("exit with error", slog.Int("code", exit_err.ExitCode()), slog.String("err", err.Error()))

		fmt.Println(err.Error())
		os.Exit(exit_err.ExitCode())
	},
}

package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"khepri.dev/horus/cmd/conf"
)

var Commands = []*cli.Command{
	{
		Name:        "init",
		Description: "initialize DB",
		Action: func(ctx_ *cli.Context) error {
			ctx := ctx_.Context
			c := conf.From(ctx)

			if c.Hr.Connect.With != "db" {
				return fmt.Errorf(`client does not connected to DB; use the "--force" to temporarily override the config to connect to the DB`)
			}

			db, err := c.Hr.Connect.Db.Open()
			if err != nil {
				return fmt.Errorf("open DB: %w", err)
			}
			if err := db.Schema.Create(ctx); err != nil {
				return fmt.Errorf("init DB: %w", err)
			}

			return nil
		},
	},
	CmdGet,
	CmdSet,
	CmdCreate,
	// CmdDelete,
	conf.CmdVersion,
}

package cmd

import "github.com/urfave/cli/v2"

var Commands = []*cli.Command{
	{
		Name:        "init",
		Description: "initialize DB",
		Action: func(ctx *cli.Context) error {
			c := ConfFrom(ctx.Context)
			c.Client.Db.WithInit = true

			_, err := c.Client.connect(ctx.Context)
			if err != nil {
				return err
			}

			return nil
		},
	},
	CmdCreate,
	CmdGet,
	// CmdDelete,
	CmdInvite,
	CmdAccept,
}

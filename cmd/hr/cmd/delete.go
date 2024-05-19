package cmd

import "github.com/urfave/cli/v2"

var CmdDelete = &cli.Command{
	Name: "delete",
	Subcommands: []*cli.Command{
		{
			Name: "all",
			Subcommands: []*cli.Command{
				CmdDeleteAllAccessTokens,
			},
		},
	},
}

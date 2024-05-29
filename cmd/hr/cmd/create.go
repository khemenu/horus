package cmd

import "github.com/urfave/cli/v2"

var CmdCreate = &cli.Command{
	Name: "create",
	Subcommands: []*cli.Command{
		CmdCreateUser,
		CmdCreateRefreshToken,
		CmdCreateAccessToken,
		CmdCreateSilo,
		CmdCreateTeam,
	},
}

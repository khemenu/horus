package cmd

import "github.com/urfave/cli/v2"

var CmdCreate = &cli.Command{
	Name: "create",
	Subcommands: []*cli.Command{
		CmdCreateUser,
		CmdCreateAccount,
		CmdCreateRefreshToken,
		CmdCreateAccessToken,
		CmdCreateSilo,
		CmdCreateTeam,
	},
}

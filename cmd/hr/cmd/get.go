package cmd

import "github.com/urfave/cli/v2"

var CmdGet = &cli.Command{
	Name: "get",
	Subcommands: []*cli.Command{
		CmdGetConfig,
		CmdGetUser,
		CmdGetToken,
	},
}

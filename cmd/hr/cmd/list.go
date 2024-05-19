package cmd

import "github.com/urfave/cli/v2"

var CmdList = &cli.Command{
	Name:        "list",
	Subcommands: []*cli.Command{},
}

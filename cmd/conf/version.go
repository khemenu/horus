package conf

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	VcsRef  = "local"
	Version = "v0.0.1"
)

var CmdVersion = &cli.Command{
	Name: "version",
	Action: func(ctx *cli.Context) error {
		fmt.Println(Version, VcsRef)
		return nil
	},
}

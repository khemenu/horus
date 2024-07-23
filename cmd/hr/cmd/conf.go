package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
)

var CmdSetConfig = &cli.Command{
	Name: "config",
	Args: true,

	Aliases:   []string{"conf"},
	ArgsUsage: "<KEY> <VALUE>",
	// ArgsUsage: "<KEY> <TYPE:VALUE> [at <PATH>]",

	Action: func(ctx *cli.Context) error {
		switch ctx.Args().Len() {
		case 2:
		case 4:
			if p := ctx.Args().Get(3); p != "at" {
				return fmt.Errorf(`expected a preposition "at" but found %s`, p)
			}

		default:
			return fmt.Errorf("requires exactly 2 or 4 arguments")
		}

		var (
			key  = ctx.Args().Get(0)
			data = ctx.Args().Get(1)
			// path = ctx.Args().Get(3)
		)

		conf := ConfFrom(ctx.Context)
		c, err := conf.Client.connect(ctx.Context)
		if err != nil {
			return err
		}

		_, err = c.Conf().Create(ctx.Context, &horus.CreateConfRequest{
			Id:    key,
			Value: data,
		})
		if err != nil {
			return err
		}

		return nil
	},
}

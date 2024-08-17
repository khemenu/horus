package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/hr/env"
)

var CmdSetConfig = &cli.Command{
	Name: "config",
	Args: true,

	Aliases:   []string{"conf"},
	ArgsUsage: "<KEY> <VALUE>",
	// ArgsUsage: "<KEY> <TYPE:VALUE> [at <PATH>]",

	Action: func(ctx_ *cli.Context) error {
		ctx := ctx_.Context
		args := ctx_.Args()
		switch args.Len() {
		case 2:

		default:
			return fmt.Errorf("requires exactly 2 arguments")
		}

		var (
			key = args.Get(0)
			val = args.Get(1)
		)

		e := env.From(ctx)
		h, err := e.Connect(ctx)
		if err != nil {
			return err
		}

		req := &horus.CreateConfRequest{Id: key, Value: val}
		_, err = h.Conf().Create(ctx, req)
		if err != nil {
			return err
		}

		return nil
	},
}

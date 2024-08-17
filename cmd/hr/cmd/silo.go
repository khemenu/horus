package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/hr/env"
)

var CmdCreateSilo = &cli.Command{
	Name:      "silo",
	Args:      true,
	ArgsUsage: " [SILO_ALIAS]",
	Action: func(ctx_ *cli.Context) error {
		ctx := ctx_.Context
		args := ctx_.Args()
		switch args.Len() {
		case 0:
		case 1:

		default:
			return fmt.Errorf("requires either no arguments or exactly 1 argument")
		}

		var (
			sil_alias = args.Get(0)
		)

		e := env.From(ctx)
		h, err := e.Connect(ctx)
		if err != nil {
			return err
		}

		req := &horus.CreateSiloRequest{}
		if sil_alias != "" {
			req.Alias = &sil_alias
		}

		v, err := h.Silo().Create(ctx, req)
		if err != nil {
			return fmt.Errorf("execute: %w", err)
		}

		o := uuid.UUID(v.Id).String()
		return e.Report(v, o)
	},
}

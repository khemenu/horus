package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/hr/env"
)

var CmdCreateUser = &cli.Command{
	Name:      "user",
	Args:      true,
	ArgsUsage: " [USER_ALIAS]",
	Action: func(ctx_ *cli.Context) error {
		ctx := ctx_.Context
		args := ctx_.Args()
		switch args.Len() {
		case 0:
		case 1:

		default:
			return fmt.Errorf("requires either no arguments or exactly 1 argument")
		}

		e := env.From(ctx)
		h, err := e.Connect(ctx)
		if err != nil {
			return err
		}

		user_alias := args.Get(0)

		req := &horus.CreateUserRequest{}
		if args.Len() > 0 {
			req.Alias = &user_alias
		}

		v, err := h.User().Create(ctx, req)
		if err != nil {
			return fmt.Errorf("execute: %w", err)
		}

		o := uuid.UUID(v.Id).String()
		return e.Report(v, o)
	},
}

var CmdGetUser = &cli.Command{
	Name:      "user",
	Args:      true,
	ArgsUsage: " <USER_ID>",
	Action: func(ctx_ *cli.Context) error {
		ctx := ctx_.Context
		args := ctx_.Args()
		switch args.Len() {
		case 1:

		default:
			return fmt.Errorf("requires exactly 1 argument")
		}

		e := env.From(ctx)
		h, err := e.Connect(ctx)
		if err != nil {
			return err
		}

		v, err := h.User().Get(ctx, horus.UserByQuery(args.Get(0)))
		if err != nil {
			return fmt.Errorf("execute: %w", err)
		}

		o := fmt.Sprintf("%s %s", uuid.UUID(v.Id), v.Alias)
		return e.Report(v, o)
	},
}

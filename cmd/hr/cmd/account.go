package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/hr/env"
)

var CmdCreateAccount = &cli.Command{
	Name:      "account",
	Args:      true,
	ArgsUsage: " [ACCOUNT_ALIAS] for <USER_ID> in <SILO_ID>",
	Action: func(ctx_ *cli.Context) error {
		ctx := ctx_.Context
		args := ctx_.Args()

		off := 0
		switch args.Len() {
		case 4:
			off++
		case 5:

		default:
			return fmt.Errorf("requires exactly 4 or 5 arguments")
		}
		if p := args.Get(off + 1); p != "for" {
			return fmt.Errorf(`expected a preposition "for" but found %s`, p)
		}
		if p := args.Get(off + 3); p != "in" {
			return fmt.Errorf(`expected a preposition "in" but found %s`, p)
		}

		var (
			acct_alias = args.Get(off + 0)
			user_id    = args.Get(off + 2)
			silo_id    = args.Get(off + 4)
		)

		e := env.From(ctx)
		h, err := e.Connect(ctx)
		if err != nil {
			return err
		}

		req := &horus.CreateAccountRequest{
			Owner: horus.UserByQuery(user_id),
			Silo:  horus.SiloByQuery(silo_id),
		}
		if acct_alias != "" {
			req.Alias = &acct_alias
		}

		v, err := h.Account().Create(ctx, req)
		if err != nil {
			return fmt.Errorf("execute: %w", err)
		}

		o := uuid.UUID(v.Id).String()
		return e.Report(v, o)
	},
}

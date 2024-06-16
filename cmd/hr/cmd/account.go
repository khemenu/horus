package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
)

var CmdCreateAccount = &cli.Command{
	Name:      "account",
	Args:      true,
	ArgsUsage: " [ACCOUNT_ALIAS] for <USER_ID> in <SILO_ID>",
	Action: func(ctx *cli.Context) error {
		var (
			acct_alias string
			user_id    string
			silo_id    string
		)
		if ctx.Args().Len() < 4 {
			return fmt.Errorf(`requires at least 4 arguments, e.g. "for", <USER_ID>, "in", and <SILO_ID>`)
		}
		if ctx.Args().Len() > 5 {
			return fmt.Errorf(`accepts up to 5 arguments`)
		}

		n := 0
		if ctx.Args().Len() == 5 {
			n++
			acct_alias = ctx.Args().Get(0)
		}
		if p := ctx.Args().Get(n); p != "for" {
			return fmt.Errorf(`expected a preposition "for" but found %s`, p)
		}
		if p := ctx.Args().Get(n + 2); p != "in" {
			return fmt.Errorf(`expected a preposition "in" but found %s`, p)
		}

		user_id = ctx.Args().Get(n + 1)
		silo_id = ctx.Args().Get(n + 3)

		conf := ConfFrom(ctx.Context)
		c, err := conf.Client.connect(ctx.Context)
		if err != nil {
			return err
		}

		user_by := horus.UserByAlias(user_id)
		if id, err := uuid.Parse(user_id); err == nil {
			user_by = horus.UserById(id)
		}

		silo_by := horus.SiloByAlias(silo_id)
		if id, err := uuid.Parse(silo_id); err == nil {
			silo_by = horus.SiloById(id)
		}

		v, err := c.Account().Create(ctx.Context, &horus.CreateAccountRequest{
			Alias: &acct_alias,
			Role:  horus.Role_ROLE_MEMBER,
			Owner: user_by,
			Silo:  silo_by,
		})
		if err != nil {
			return err
		}

		o := uuid.UUID(v.Id).String()
		return conf.Reporter.Report(v, o)
	},
}

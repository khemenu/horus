package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
)

var CmdInvite = &cli.Command{
	Name: "invite",
	Subcommands: []*cli.Command{
		CmdInviteUser,
	},
}

var CmdInviteUser = &cli.Command{
	Name:      "user",
	Args:      true,
	ArgsUsage: " <USER_ID> in <SILO_ID>",

	Action: func(ctx *cli.Context) error {
		if ctx.Args().Len() != 3 {
			return fmt.Errorf(`requires exact 3 arguments`)
		}
		if p := ctx.Args().Get(1); p != "in" {
			return fmt.Errorf(`expected a preposition "in" but found %s`, p)
		}

		user_id := ctx.Args().Get(0)
		silo_id := ctx.Args().Get(2)

		conf := ConfFrom(ctx.Context)
		if err := conf.Client.notToBeBareServe(); err != nil {
			return err
		}

		c, err := conf.Client.connect(ctx.Context)
		if err != nil {
			return err
		}

		silo_uuid, _ := uuid.Parse(silo_id)
		v, err := c.Invitation().Create(ctx.Context, &horus.CreateInvitationRequest{
			Invitee: user_id,
			Type:    horus.InvitationTypeInternal,

			Silo: &horus.Silo{
				Id:    silo_uuid[:],
				Alias: silo_id,
			},
		})
		if err != nil {
			return err
		}

		o := uuid.UUID(v.Id).String()
		return conf.Reporter.Report(v, o)
	},
}

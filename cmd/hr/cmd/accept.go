package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
)

var CmdAccept = &cli.Command{
	Name:      "accept",
	Args:      true,
	ArgsUsage: " <INVITATION_UUID>",

	Action: func(ctx *cli.Context) error {
		if ctx.Args().Len() != 1 {
			return fmt.Errorf(`requires exact 1 argument`)
		}

		invitation_uuid, err := uuid.Parse(ctx.Args().Get(0))
		if err != nil {
			return fmt.Errorf("invalid UUID: %w", err)
		}

		conf := ConfFrom(ctx.Context)
		if err := conf.Client.notToBeBareServe(); err != nil {
			return err
		}

		c, err := conf.Client.connect(ctx.Context)
		if err != nil {
			return err
		}

		_, err = c.Invitation().Accept(ctx.Context, &horus.AcceptInvitationRequest{
			Id: invitation_uuid[:],
		})
		if err != nil {
			return err
		}

		return nil
	},
}

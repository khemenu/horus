package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
)

var CmdCreateUser = &cli.Command{
	Name:      "user",
	Args:      true,
	ArgsUsage: " USERNAME",
	Action: func(ctx *cli.Context) error {
		c := ConfFrom(ctx.Context)
		s := c.Client.mustConnect(ctx.Context)

		if ctx.Args().Len() == 0 {
			return fmt.Errorf("USERNAME not given")
		}

		name := ctx.Args().Get(0)
		v, err := s.User().Create(ctx.Context, &horus.CreateUserRequest{
			User: &horus.User{
				Name: name,
			},
		})
		if err != nil {
			return err
		}

		fmt.Println(uuid.UUID(v.Id).String())
		return nil
	},
}

var CmdUserGet = &cli.Command{
	Name:      "user",
	Args:      true,
	ArgsUsage: " USER_ID",
	Action: func(ctx *cli.Context) error {
		c := ConfFrom(ctx.Context)
		s := c.Client.mustConnect(ctx.Context)

		if ctx.Args().Len() == 0 {
			return fmt.Errorf("USER_ID not given")
		}

		user_id, err := uuid.Parse(ctx.Args().Get(0))
		if err != nil {
			return fmt.Errorf("USER_ID must be UUID: %w", err)
		}

		v, err := s.User().Get(ctx.Context, &horus.GetUserRequest{
			Id:   user_id[:],
			View: horus.GetUserRequest_WITH_EDGE_IDS,
		})
		if err != nil {
			return err
		}

		fmt.Print(v)
		return nil
	},
}

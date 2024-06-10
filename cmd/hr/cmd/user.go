package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
	"khepri.dev/horus/ent/user"
)

var CmdCreateUser = &cli.Command{
	Name:      "user",
	Args:      true,
	ArgsUsage: " [USER_ALIAS]",
	Action: func(ctx *cli.Context) error {
		conf := ConfFrom(ctx.Context)
		c, err := conf.Client.connect(ctx.Context)
		if err != nil {
			return err
		}

		var alias string
		if ctx.Args().Len() > 0 {
			alias = ctx.Args().Get(0)
		}

		v, err := c.User().Create(ctx.Context, &horus.CreateUserRequest{
			Alias: &alias,
		})
		if err != nil {
			return err
		}

		o := uuid.UUID(v.Id).String()
		return conf.Reporter.Report(v, o)
	},
}

var CmdGetUser = &cli.Command{
	Name:      "user",
	Args:      true,
	ArgsUsage: " <USER_ID>",
	Action: func(ctx *cli.Context) error {
		if ctx.Args().Len() == 0 {
			return fmt.Errorf("<USER_ID> is required")
		}

		conf := ConfFrom(ctx.Context)
		c, err := conf.Client.connectDbServer(ctx.Context)
		if err != nil {
			return err
		}

		user_id := ctx.Args().Get(0)
		user_uuid, err := uuid.Parse(ctx.Args().Get(0))
		if err == nil {
			goto Q
		}
		if !conf.Client.isBareServer() {
			return fmt.Errorf("USER_ID must be valid UUID: %w", err)
		} else if user_uuid, err = conf.Client.db.
			User.Query().
			Where(user.AliasEQ(user_id)).
			OnlyID(ctx.Context); err != nil {
			return fmt.Errorf("query user by alias: %w", err)
		}

	Q:
		v, err := c.User().Get(ctx.Context, &horus.GetUserRequest{Key: &horus.GetUserRequest_Id{
			Id: user_uuid[:],
		}})
		if err != nil {
			return err
		}

		o := fmt.Sprintf("%s %s", uuid.UUID(v.Id), v.Alias)
		return conf.Reporter.Report(v, o)
	},
}

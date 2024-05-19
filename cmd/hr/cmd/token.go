package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
	"khepri.dev/horus/ent/predicate"
	"khepri.dev/horus/ent/token"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/service/frame"
)

func actCreateToken(ctx *cli.Context, token_type string) error {
	c := ConfFrom(ctx.Context)
	s := c.Client.mustBareConnect(ctx.Context)

	if ctx.Args().Len() == 0 {
		return fmt.Errorf("USER_ID not given")
	}

	var pred predicate.User
	if user_id, err := uuid.Parse(ctx.Args().Get(0)); err == nil {
		pred = user.IDEQ(user_id)
	} else {
		pred = user.NameEQ(ctx.Args().Get(0))
	}

	user, err := s.client.User.Query().Where(pred).Only(ctx.Context)
	if err != nil {
		return err
	}

	ctx.Context = frame.WithContext(ctx.Context, &frame.Frame{
		Actor: user,
	})
	token, err := s.svc.Token().Create(ctx.Context, &horus.CreateTokenRequest{Token: &horus.Token{
		Type: token_type,
	}})
	if err != nil {
		return err
	}

	fmt.Println(token.Value)
	return nil
}

var CmdCreateRefreshToken = &cli.Command{
	Name: "refresh-token",
	Subcommands: []*cli.Command{
		{
			Name:      "for",
			Args:      true,
			ArgsUsage: " USER_ID",
			Action: func(ctx *cli.Context) error {
				return actCreateToken(ctx, horus.TokenTypeRefresh)
			},
		},
	},
}

var CmdCreateAccessToken = &cli.Command{
	Name: "access-token",
	Subcommands: []*cli.Command{
		{
			Name:      "for",
			Args:      true,
			ArgsUsage: " USER_ID",
			Action: func(ctx *cli.Context) error {
				return actCreateToken(ctx, horus.TokenTypeAccess)
			},
		},
	},
}

var CmdDeleteAllAccessTokens = &cli.Command{
	Name: "access-tokens",
	Subcommands: []*cli.Command{
		{
			Name:      "of",
			Args:      true,
			ArgsUsage: " USER_ID",
			Action: func(ctx *cli.Context) error {
				c := ConfFrom(ctx.Context)
				s := c.Client.mustBareConnect(ctx.Context)

				if ctx.Args().Len() == 0 {
					return fmt.Errorf("USER_ID not given")
				}

				var pred predicate.User
				if user_id, err := uuid.Parse(ctx.Args().Get(0)); err == nil {
					pred = user.IDEQ(user_id)
				} else {
					pred = user.NameEQ(ctx.Args().Get(0))
				}

				owner, err := s.client.User.Query().Where(pred).Only(ctx.Context)
				if err != nil {
					return err
				}

				n, err := s.client.Token.Delete().
					Where(token.And(
						token.TypeEQ(horus.TokenTypeAccess),
						token.HasOwnerWith(user.IDEQ(owner.ID)),
					)).
					Exec(ctx.Context)
				if err != nil {
					return err
				}

				fmt.Println(n)
				return nil
			},
		},
	},
}

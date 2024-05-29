package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
)

func actCreateBearerToken(ctx *cli.Context, token_type string) error {
	conf := ConfFrom(ctx.Context)
	c, err := conf.Client.connect(ctx.Context)
	if err != nil {
		return err
	}

	v, err := c.Token().Create(ctx.Context, &horus.CreateTokenRequest{Token: &horus.Token{
		Type: token_type,
	}})
	if err != nil {
		return err
	}

	o := v.Value
	return conf.Reporter.Report(v, o)
}

var CmdCreateRefreshToken = &cli.Command{
	Name: "refresh-token",
	Action: func(ctx *cli.Context) error {
		return actCreateBearerToken(ctx, horus.TokenTypeRefresh)
	},
}

var CmdCreateAccessToken = &cli.Command{
	Name: "access-token",
	Action: func(ctx *cli.Context) error {
		return actCreateBearerToken(ctx, horus.TokenTypeAccess)
	},
}

var CmdGetToken = &cli.Command{
	Name:      "token",
	Args:      true,
	ArgsUsage: " TOKEN_UUID",
	Action: func(ctx *cli.Context) error {
		conf := ConfFrom(ctx.Context)
		c, err := conf.Client.connect(ctx.Context)
		if err != nil {
			return err
		}

		if ctx.Args().Len() == 0 {
			return fmt.Errorf("TOKEN_UUID must be provided")
		}

		token_uuid, err := uuid.Parse(ctx.Args().Get(0))
		if err != nil {
			return fmt.Errorf("invalid UUID")
		}

		v, err := c.Token().Get(ctx.Context, &horus.GetTokenRequest{
			Id: token_uuid[:],
		})
		if err != nil {
			return err
		}

		o := fmt.Sprintf("%s valid until %s", v.Type, v.DateExpired.AsTime().String())
		return conf.Reporter.Report(v, o)
	},
}

// var CmdDeleteAllAccessTokens = &cli.Command{
// 	Name: "access-tokens",
// 	Subcommands: []*cli.Command{
// 		{
// 			Name:      "of",
// 			Args:      true,
// 			ArgsUsage: " USER_ID",
// 			Action: func(ctx *cli.Context) error {
// 				conf := ConfFrom(ctx.Context)
// 				s := conf.Client.mustBareServer(ctx.Context)

// 				if ctx.Args().Len() == 0 {
// 					return fmt.Errorf("USER_ID not given")
// 				}

// 				var pred predicate.User
// 				if user_id, err := uuid.Parse(ctx.Args().Get(0)); err == nil {
// 					pred = user.IDEQ(user_id)
// 				} else {
// 					pred = user.AliasEQ(ctx.Args().Get(0))
// 				}

// 				owner, err := s.db.User.Query().Where(pred).Only(ctx.Context)
// 				if err != nil {
// 					return err
// 				}

// 				n, err := s.db.Token.Delete().
// 					Where(token.And(
// 						token.TypeEQ(horus.TokenTypeAccess),
// 						token.HasOwnerWith(user.IDEQ(owner.ID)),
// 					)).
// 					Exec(ctx.Context)
// 				if err != nil {
// 					return err
// 				}

// 				fmt.Println(n)
// 				return nil
// 			},
// 		},
// 	},
// }

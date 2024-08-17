package cmd

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/hr/env"
)

func actCreateBearerToken(ctx context.Context, token_type string) error {
	e := env.From(ctx)
	if err := e.NotToBeBareServe(); err != nil {
		// Token value should be encrypted.
		return err
	}

	h, err := e.Connect(ctx)
	if err != nil {
		return err
	}

	v, err := h.Token().Create(ctx, &horus.CreateTokenRequest{
		Type: token_type,
	})
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	o := v.Value
	return e.Report(v, o)
}

var CmdCreateRefreshToken = &cli.Command{
	Name: "refresh-token",
	Action: func(ctx *cli.Context) error {
		return actCreateBearerToken(ctx.Context, horus.TokenTypeRefresh)
	},
}

var CmdCreateAccessToken = &cli.Command{
	Name: "access-token",
	Action: func(ctx *cli.Context) error {
		return actCreateBearerToken(ctx.Context, horus.TokenTypeAccess)
	},
}

var CmdGetToken = &cli.Command{
	Name:      "token",
	Args:      true,
	ArgsUsage: " TOKEN_UUID",
	Action: func(ctx_ *cli.Context) error {
		ctx := ctx_.Context
		args := ctx_.Args()
		switch args.Len() {
		case 1:

		default:
			return fmt.Errorf("requires exactly 1 argument")
		}

		var (
			token_uuid uuid.UUID
			err        error
		)
		if token_uuid, err = uuid.Parse(args.Get(0)); err != nil {
			return fmt.Errorf("invalid UUID")
		}

		e := env.From(ctx)
		h, err := e.Connect(ctx)
		if err != nil {
			return err
		}

		req := horus.TokenById(token_uuid)
		v, err := h.Token().Get(ctx, req)
		if err != nil {
			return fmt.Errorf("execute: %w", err)
		}

		o := fmt.Sprintf("%s %s", uuid.UUID(v.Id), v.DateExpired)
		return e.Report(v, o)
	},
}

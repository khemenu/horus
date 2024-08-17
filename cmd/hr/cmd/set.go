package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/hr/env"
)

var CmdSet = &cli.Command{
	Name: "set",
	Subcommands: []*cli.Command{
		CmdSetPassword,
	},
}

var CmdSetPassword = &cli.Command{
	Name:    "password",
	Aliases: []string{"pw"},

	Action: func(ctx_ *cli.Context) error {
		ctx := ctx_.Context
		if ctx_.Args().Len() > 0 {
			return fmt.Errorf("no arguments are required")
		}

		e := env.From(ctx)
		if err := e.NotToBeBareServe(); err != nil {
			return err
		}

		var (
			pw  string
			err error
		)
		if !term.IsTerminal(syscall.Stdin) {
			pw, err = bufio.NewReader(os.Stdin).ReadString('\n')
		} else {
			var pw_ []byte
			fmt.Print("password: ")
			pw_, err = term.ReadPassword(syscall.Stdin)
			pw = string(pw_)
		}
		if err != nil {
			return fmt.Errorf("read password from stdin: %w", err)
		}
		pw = strings.TrimSpace(pw)

		h, err := e.Connect(ctx)
		if err != nil {
			return err
		}

		req := &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		}
		_, err = h.Token().Create(ctx, req)
		if err != nil {
			return err
		}

		return nil
	},
}

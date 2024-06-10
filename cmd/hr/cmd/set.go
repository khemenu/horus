package cmd

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
	"khepri.dev/horus"
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

	Action: func(ctx *cli.Context) error {
		conf := ConfFrom(ctx.Context)
		if err := conf.Client.notToBeBareServe(); err != nil {
			return err
		}

		c, err := conf.Client.connect(ctx.Context)
		if err != nil {
			return err
		}

		var pw string
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

		_, err = c.Token().Create(ctx.Context, &horus.CreateTokenRequest{
			Value: string(pw),
			Type:  horus.TokenTypeBasic,
		})
		if err != nil {
			return err
		}

		return nil
	},
}

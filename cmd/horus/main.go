package main

import (
	"context"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"khepri.dev/horus/cmd/horus/cmd"
)

func main() {
	c, err := cmd.ParseArgs(os.Args)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err := cmd.Run(ctx, c); err != nil {
		c.Log.NewLogger().Error("fatal", "err", err)
		os.Exit(1)
	}
}

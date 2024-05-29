package main

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
	"khepri.dev/horus/cmd/hr/cmd"
)

func main() {
	cmd.App.Run(os.Args)
}

package main

import (
	"fmt"
	"os"

	"github.com/ipfs-force-community/api-compare/cmd"
	"github.com/ipfs-force-community/api-compare/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Compare the apis of venus and lotus",
		Usage: "",
		Flags: []cli.Flag{
			cmd.LotusURLFlag,
			cmd.LotusTokenFlag,
			cmd.VenusURLFlag,
			cmd.VenusTokenFlag,
			&cli.IntFlag{
				Name:  "start-height",
				Usage: "Start comparing the height of the API",
			},
			&cli.IntFlag{
				Name:  "stop-height",
				Usage: "Stop after running n heights",
			},
			&cli.IntFlag{
				Name:  "concurrency",
				Value: 1,
			},
			&cli.BoolFlag{
				Name:  "enable-eth-rpc",
				Usage: "Need to compare ETH interfaces",
			},
		},
		Version: version.UserVersion(),
		Action:  cmd.Run,
	}

	app.Setup()

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
	}
}

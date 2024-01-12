package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ipfs-force-community/api-compare/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Compare the apis of venus and lotus",
		Usage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "url",
				Value: "ws://192.168.200.18:9944",
			},
			&cli.StringFlag{
				Name: "token",
			},
		},
		Version: version.UserVersion(),
		Action:  run,
	}

	app.Setup()

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
	}
}

func run(cctx *cli.Context) error {
	url := cctx.String("url")
	token := cctx.String("token")
	fmt.Println("url:", url, "token:", token)

	ctx, cancel := context.WithCancel(cctx.Context)
	defer cancel()

	api, close, err := dialSubspaceRPC(ctx, url, token, nil)
	if err != nil {
		return err
	}
	defer close()

	solutionRangeChan, err := api.SubscribeSlotInfo(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case solutionRange := <-solutionRangeChan:
			fmt.Println(solutionRange)
		}
	}
}

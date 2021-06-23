package main

import (
	"log"
	"os"

	"github.com/eventflowdb/eventflowdb/flowctl"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "flowctl",
		Commands: []*cli.Command{
			flowctl.VersionCommand,
			flowctl.SizeCommand,
			flowctl.GetCommand,
			flowctl.GetAllCommand,
			flowctl.EventCountCommand,
			flowctl.EventCountEstimateCommand,
			flowctl.StreamCountCommand,
			flowctl.StreamCountEstimateCommand,
			flowctl.ListStreamsCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

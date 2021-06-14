package main

import (
	"log"
	"os"

	"github.com/kajjagtenberg/eventflowdb/flowctl"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "flowctl",
		Commands: []*cli.Command{
			flowctl.VersionCommand,
			flowctl.SizeCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

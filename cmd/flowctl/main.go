package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	Version = "0.1.0"
)

func main() {
	app := &cli.App{
		Name: "flowctl",
		Commands: []cli.Command{
			{
				Name:      "version",
				ShortName: "v",
				Usage:     "Returns the cli version",
				Action: func(c *cli.Context) error {
					log.Println(Version)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

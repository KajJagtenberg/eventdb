package flowctl

import (
	"log"

	"github.com/urfave/cli/v2"
)

const (
	Version = "0.1.0"
)

var VersionCommand = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "Returns the cli version",
	Action: func(c *cli.Context) error {
		log.Println(Version)
		return nil
	},
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/kajjagtenberg/eventflowdb/api"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

const (
	Version = "0.1.0"
)

func main() {
	app := &cli.App{
		Name: "flowctl",
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Returns the cli version",
				Action: func(c *cli.Context) error {
					log.Println(Version)
					return nil
				},
			},
			{
				Name:    "size",
				Aliases: []string{"s"},
				Usage:   "Returns the storage size",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "address",
						Usage:   "Address of the cluster",
						EnvVars: []string{"ADDRESS"},
						Value:   "127.0.0.1:6543",
					},
					&cli.StringFlag{
						Name:    "port",
						Usage:   "Port of the cluster",
						EnvVars: []string{"PORT"},
						Value:   "6543",
					},
					&cli.BoolFlag{
						Name:  "human",
						Usage: "Whether to report the size in a human readable format",
					},
				},
				Action: func(c *cli.Context) error {
					address := c.String("address")
					human := c.Bool("human")

					conn, err := grpc.Dial(address, grpc.WithInsecure())
					if err != nil {
						return err
					}
					defer conn.Close()

					store := api.NewEventStoreServiceClient(conn)
					res, err := store.Size(context.Background(), &api.SizeRequest{})
					if err != nil {
						return err
					}

					if human {
						log.Println(res.SizeHuman)
					} else {
						log.Println(res.Size)
					}

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

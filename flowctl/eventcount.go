package flowctl

import (
	"context"
	"log"

	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var EventCountCommand = &cli.Command{
	Name:    "eventcount",
	Aliases: []string{"ec"},
	Usage:   "Returns the total event count",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Address of the cluster",
			EnvVars: []string{"ADDRESS"},
			Value:   "127.0.0.1:6543",
		},
	},
	Action: func(c *cli.Context) error {
		address := c.String("address")

		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()

		store := api.NewEventStoreServiceClient(conn)
		res, err := store.EventCount(context.Background(), &api.EventCountRequest{})
		if err != nil {
			return err
		}

		log.Println(res.Count)

		return nil
	},
}

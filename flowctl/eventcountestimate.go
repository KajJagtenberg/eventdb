package flowctl

import (
	"context"
	"log"

	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var EventCountEstimateCommand = &cli.Command{
	Name:    "eventcountestimate",
	Aliases: []string{"ece"},
	Usage:   "Returns the total event count from the cache",
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

		store := api.NewEventStoreClient(conn)
		res, err := store.EventCountEstimate(context.Background(), &api.EventCountEstimateRequest{})
		if err != nil {
			return err
		}

		log.Println(res.Count)

		return nil
	},
}

package flowctl

import (
	"context"
	"log"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var StreamCountCommand = &cli.Command{
	Name:    "streamcount",
	Aliases: []string{"sc"},
	Usage:   "Returns the total stream count",
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
		res, err := store.StreamCount(context.Background(), &api.StreamCountRequest{})
		if err != nil {
			return err
		}

		log.Println(res.Count)

		return nil
	},
}

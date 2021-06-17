package flowctl

import (
	"context"
	"encoding/json"
	"os"

	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var ListStreamsCommand = &cli.Command{
	Name:    "liststream",
	Aliases: []string{"ls"},
	Usage:   "Returns a list of streams",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Address of the cluster",
			EnvVars: []string{"ADDRESS"},
			Value:   "127.0.0.1:6543",
		},
		&cli.IntFlag{
			Name:    "skip",
			Aliases: []string{"s"},
			Usage:   "The amount of streams to skip in the results",
		},
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"v"},
			Usage:   "The max number of streams to return",
			Value:   0,
		},
	},
	Action: func(c *cli.Context) error {
		address := c.String("address")
		skip := c.Int("skip")
		limit := c.Int("limit")

		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()

		store := api.NewEventStoreClient(conn)
		res, err := store.ListStreams(context.Background(), &api.ListStreamsRequest{
			Skip:  uint32(skip),
			Limit: uint32(limit),
		})
		if err != nil {
			return err
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(res.Streams)

		return nil
	},
}

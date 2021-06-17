package flowctl

import (
	"context"
	"log"

	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var GetCommand = &cli.Command{
	Name:    "get",
	Aliases: []string{"g"},
	Usage:   "Returns events for given stream",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Address of the cluster",
			EnvVars: []string{"ADDRESS"},
			Value:   "127.0.0.1:6543",
		},
		&cli.StringFlag{
			Name:    "stream",
			Aliases: []string{"s"},
			Usage:   "The given stream",
		},
		&cli.IntFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "The offset from the start of the stream",
			Value:   0,
		},
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"l"},
			Usage:   "The maximum amount of events to return",
		},
	},
	Action: func(c *cli.Context) error {
		address := c.String("address")
		stream := c.String("stream")
		version := c.Int("version")
		limit := c.Int("limit")

		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()

		store := api.NewEventStoreClient(conn)
		res, err := store.Get(context.Background(), &api.GetRequest{
			Stream:  stream,
			Version: uint32(version),
			Limit:   uint32(limit),
		})
		if err != nil {
			return err
		}

		log.Println(res.Events)

		return nil
	},
}

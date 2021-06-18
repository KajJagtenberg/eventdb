package flowctl

import (
	"context"
	"encoding/json"
	"log"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var GetAllCommand = &cli.Command{
	Name:    "getall",
	Aliases: []string{"ga"},
	Usage:   "Returns events from the global stream",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Address of the cluster",
			EnvVars: []string{"ADDRESS"},
			Value:   "127.0.0.1:6543",
		},
		&cli.StringFlag{
			Name:    "offset",
			Aliases: []string{"s"},
			Usage:   "The offset in the global stream",
			Value:   "00000000000000000000000000",
		},
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"l"},
			Usage:   "The maximum amount of events to return",
		},
	},
	Action: func(c *cli.Context) error {
		address := c.String("address")
		offset := c.String("offset")
		limit := c.Int("limit")

		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()

		store := api.NewEventStoreClient(conn)
		res, err := store.GetAll(context.Background(), &api.GetAllRequest{
			Offset: offset,
			Limit:  uint32(limit),
		})
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(res.Events, "", "  ")
		if err != nil {
			return err
		}

		log.Println(string(data))

		return nil
	},
}

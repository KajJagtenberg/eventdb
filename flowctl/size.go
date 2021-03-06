package flowctl

import (
	"context"
	"log"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var SizeCommand = &cli.Command{
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

		store := api.NewEventStoreClient(conn)
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
}

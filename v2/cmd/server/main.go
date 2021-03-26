package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/helmet/v2"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/cluster"
	"github.com/kajjagtenberg/eventflowdb/env"
	"go.etcd.io/bbolt"
)

var (
	localID       = env.GetEnv("RAFT_LOCAL_ID", "")
	bindAddr      = env.GetEnv("RAFT_BIND_ADDR", ":6542")
	advrAddr      = env.GetEnv("RAFT_ADVR_ADDR", bindAddr)
	bootstrap     = env.GetEnv("RAFT_BOOTSTRAP", "false") == "true"
	stateLocation = env.GetEnv("STATE_LOCATION", "data/state.dat")
)

func main() {
	db, err := bbolt.Open(stateLocation, 0666, bbolt.DefaultOptions)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	fsm, err := cluster.NewFSM(db)
	if err != nil {
		log.Fatalf("Failed to create FSM: %v", err)
	}

	raft, err := cluster.NewRaftServer(localID, bindAddr, advrAddr, fsm, true)
	if err != nil {
		log.Fatalf("Failed to create Raft: %v", err)
	}
	defer raft.Shutdown()

	server := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)
	server.Use(helmet.New())
	server.Post("/add", func(c *fiber.Ctx) error {
		var body struct {
			Stream  uuid.UUID                   `json:"stream"`
			Version uint32                      `json:"version"`
			Events  []*cluster.AddCommand_Event `json:"events"`
		}

		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if len(body.Events) == 0 {
			return errors.New("Events cannot be empty")
		}

		for i, event := range body.Events {
			if len(event.Type) == 0 {
				return fmt.Errorf("Type of event at index %d cannot be empty", i)
			}

			// Validate other fields
		}

		cmd, err := proto.Marshal(&cluster.ApplyLog{
			Command: &cluster.ApplyLog_Add{
				Add: &cluster.AddCommand{
					Stream:  body.Stream[:],
					Version: body.Version,
					Events:  body.Events,
				},
			},
		})
		if err != nil {
			return err
		}

		future := raft.Apply(cmd, time.Second*5)

		if err := future.Error(); err != nil {
			return err
		}

		return c.SendString("Set")
	})

	if err := server.Listen(":3000"); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}

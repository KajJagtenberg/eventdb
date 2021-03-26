package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/helmet/v2"
	"github.com/kajjagtenberg/eventflowdb/cluster"
	"github.com/kajjagtenberg/eventflowdb/env"
)

var (
	localID   = env.GetEnv("RAFT_LOCAL_ID", "")
	bindAddr  = env.GetEnv("RAFT_BIND_ADDR", ":6542")
	advrAddr  = env.GetEnv("RAFT_ADVR_ADDR", bindAddr)
	bootstrap = env.GetEnv("RAFT_BOOTSTRAP", "false") == "true"
)

func main() {
	fsm, err := cluster.NewFSM()
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
	server.Post("/set", func(c *fiber.Ctx) error {
		var body struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}

		if err := c.BodyParser(&body); err != nil {
			return err
		}

		cmd := []byte(fmt.Sprintf("SET %s=%s", body.Key, body.Value))

		future := raft.Apply(cmd, time.Second*5)

		if err := future.Error(); err != nil {
			return err
		}

		return c.SendString("Set")
	})
	server.Get("/state", func(c *fiber.Ctx) error {
		return c.JSON(fsm.GetState())
	})

	if err := server.Listen(":3000"); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}

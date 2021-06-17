package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kajjagtenberg/eventflowdb/api"
	"github.com/kajjagtenberg/eventflowdb/example/account"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "127.0.0.1:6543", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	eventstore := api.NewEventStoreServiceClient(conn)

	app := fiber.New()
	app.Post("/api/v1/account/register", func(c *fiber.Ctx) error {
		var body struct {
			ID   uuid.UUID `json:"id"`
			Name string    `json:"name"`
		}
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		return account.Handle(eventstore, body.ID, &account.RegisterAccount{
			Name: body.Name,
		})
	})
	app.Post("/api/v1/account/changename", func(c *fiber.Ctx) error {
		var body struct {
			ID   uuid.UUID `json:"id"`
			Name string    `json:"name"`
		}
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		return account.Handle(eventstore, body.ID, &account.ChangeAccountName{
			Name: body.Name,
		})
	})
	app.Listen(":8080")
}

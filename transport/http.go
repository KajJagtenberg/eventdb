package transport

import (
	"encoding/json"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/constants"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/gofiber/fiber/v2"
)

func HTTPHandler(eventstore store.EventStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd struct {
			Operation string          `json:"op"`
			Payload   json.RawMessage `json:"payload"`
		}
		if err := c.BodyParser(&cmd); err != nil {
			return err
		}

		switch cmd.Operation {
		case "version":
			return version(c, cmd.Payload)
		}

		return fiber.NewError(fiber.StatusBadRequest, "Unknown command")
	}
}

func version(c *fiber.Ctx, payload json.RawMessage) error {
	var req api.VersionRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		return fiber.ErrUnprocessableEntity
	}

	return c.JSON(api.VersionResponse{
		Version: constants.Version,
	})
}

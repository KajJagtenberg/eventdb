package transport

import (
	"encoding/json"
	"time"

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
		case "version", "v":
			var req api.VersionRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			return c.JSON(api.VersionResponse{
				Version: constants.Version,
			})
		case "add", "a":
			var req api.AddRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.Add(&req)
			if err != nil {
				return err
			}

			return c.JSON(res)
		case "get", "g":
			var req api.GetRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.Get(&req)
			if err != nil {
				return err
			}

			return c.JSON(res)
		case "getall", "ga":
			var req api.GetAllRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.GetAll(&req)
			if err != nil {
				return err
			}

			return c.JSON(res)
		case "eventcount", "ec":
			var req api.EventCountRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.EventCount(&req)
			if err != nil {
				return err
			}

			return c.JSON(res)
		case "eventcountestimate", "ece":
			var req api.EventCountEstimateRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.EventCountEstimate(&req)
			if err != nil {
				return err
			}

			return c.JSON(res)
		case "streamcount", "sc":
			var req api.StreamCountRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.StreamCount(&req)
			if err != nil {
				return err
			}

			return c.JSON(res)
		case "streamcountestimate", "sce":
			var req api.StreamCountEstimateRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.StreamCountEstimate(&req)
			if err != nil {
				return err
			}

			return c.JSON(res)
		case "liststreams", "ls":
			var req api.ListStreamsRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.ListStreams(&req)
			if err != nil {
				return err
			}

			return c.JSON(res)
		case "size", "s":
			var req api.SizeRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.Size(&req)
			if err != nil {
				return err
			}

			return c.JSON(res)
		case "uptime", "up":
			var req api.UptimeRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			uptime := time.Since(start)

			return c.JSON(api.UptimeResponse{
				Uptime:      uptime.Milliseconds(),
				UptimeHuman: uptime.String(),
			})
		}

		return fiber.NewError(fiber.StatusBadRequest, "Unknown command")
	}
}

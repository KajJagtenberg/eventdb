package transport

import (
	"strconv"
	"time"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/constants"
	"github.com/eventflowdb/eventflowdb/env"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/sirupsen/logrus"
)

func version() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(api.VersionResponse{
			Version: constants.Version,
		})
	}
}

func stream(eventstore store.EventStore, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req api.GetStreamRequest
		req.Stream = c.Params("id")

		version, err := strconv.ParseUint(c.Query("version", "0"), 10, 32)
		if err != nil {
			return fiber.ErrBadRequest
		}
		limit, err := strconv.ParseUint(c.Query("limit", "0"), 10, 32)
		if err != nil {
			return fiber.ErrBadRequest
		}

		req.Version = uint32(version)
		req.Limit = uint32(limit)

		res, err := eventstore.GetStream(&req)
		if err != nil {
			switch err {
			default:
				logger.Println(err)
				return fiber.ErrInternalServerError
			}
		}

		return c.JSON(res)
	}
}

func all(eventstore store.EventStore, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req api.GetGlobalStreamRequest
		offset, err := strconv.ParseUint(c.Query("offset", "0"), 10, 64)
		if err != nil {
			return fiber.ErrBadRequest
		}
		limit, err := strconv.ParseUint(c.Query("limit", "0"), 10, 32)
		if err != nil {
			return fiber.ErrBadRequest
		}

		req.Offset = offset
		req.Limit = uint32(limit)

		res, err := eventstore.GetGlobalStream(&req)
		if err != nil {
			switch err {
			default:
				logger.Println(err)
				return fiber.ErrInternalServerError
			}
		}

		return c.JSON(res)
	}
}

func appendStream(eventstore store.EventStore, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req api.AppendStreamRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		res, err := eventstore.AppendStream(&req)
		if err != nil {
			switch err {
			default:
				logger.Println(err)
				return fiber.ErrInternalServerError
			}
		}

		return c.JSON(res)
	}
}

func streamCount(eventstore store.EventStore, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var res interface{}
		var err error

		if c.Query("estimate", "false") == "true" {
			res, err = eventstore.StreamCountEstimate(&api.StreamCountEstimateRequest{})
		} else {
			res, err = eventstore.StreamCount(&api.StreamCountRequest{})
		}

		if err != nil {
			switch err {
			default:
				logger.Println(err)
				return fiber.ErrInternalServerError
			}
		}

		return c.JSON(res)
	}
}

func event(eventstore store.EventStore, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req api.GetEventRequest
		req.Id = c.Params("id")

		res, err := eventstore.GetEvent(&req)
		if err != nil {
			switch err {
			default:
				logger.Println(err)
				return fiber.ErrInternalServerError
			}
		}

		return c.JSON(res)
	}
}

func eventCount(eventstore store.EventStore, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var res interface{}
		var err error

		if c.Query("estimate", "false") == "true" {
			res, err = eventstore.EventCountEstimate(&api.EventCountEstimateRequest{})
		} else {
			res, err = eventstore.EventCount(&api.EventCountRequest{})
		}

		if err != nil {
			switch err {
			default:
				logger.Println(err)
				return fiber.ErrInternalServerError
			}
		}

		return c.JSON(res)
	}
}

func uptime() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uptime := time.Since(start)

		return c.JSON(api.UptimeResponse{
			Uptime:      uptime.Milliseconds(),
			UptimeHuman: uptime.String(),
		})
	}
}

func streams(eventstore store.EventStore, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req api.ListStreamsRequest

		offset, err := strconv.ParseUint(c.Query("offset", "0"), 10, 32)
		if err != nil {
			return fiber.ErrBadRequest
		}
		limit, err := strconv.ParseUint(c.Query("limit", "0"), 10, 32)
		if err != nil {
			return fiber.ErrBadRequest
		}

		req.Skip = uint32(offset)
		req.Limit = uint32(limit)

		res, err := eventstore.ListStreams(&req)
		if err != nil {
			switch err {
			default:
				logger.Println(err)
				return fiber.ErrInternalServerError
			}
		}

		return c.JSON(res)
	}
}

func RunRestServer(eventstore store.EventStore, logger *logrus.Logger) *fiber.App {
	httpPort := env.GetEnv("HTTP_PORT", "16543")

	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	v1 := server.Group("/api/v1")
	v1.Use(compress.New())
	v1.Get("/version", version())
	v1.Get("/stream/all", all(eventstore, logger))
	v1.Get("/stream/count", streamCount(eventstore, logger))
	v1.Get("/stream/:id", stream(eventstore, logger))
	v1.Post("/stream/:id", appendStream(eventstore, logger))
	v1.Get("/event/count", eventCount(eventstore, logger))
	v1.Get("/event/:id", event(eventstore, logger))
	v1.Get("/uptime", uptime())
	v1.Get("/streams", streams(eventstore, logger))

	go func() {
		logger.Printf("REST server listening on %s", httpPort)

		if err := server.Listen(":" + httpPort); err != nil {
			logger.Fatal(err)
		}
	}()

	return server
}

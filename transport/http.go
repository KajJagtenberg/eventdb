package transport

import (
	"encoding/json"
	"time"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/constants"
	"github.com/eventflowdb/eventflowdb/env"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func HTTPHandler(eventstore store.EventStore, logger *logrus.Logger) fiber.Handler {
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
			// var req api.VersionRequest
			// if err := json.Unmarshal(cmd.Payload, &req); err != nil {
			// 	return fiber.ErrUnprocessableEntity
			// }

			return c.JSON(api.VersionResponse{
				Version: constants.Version,
			})
		case "add", "a":
			var req api.AddRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.Add(&req)
			switch err {
			case store.ErrConcurrentStreamModification, store.ErrEmptyEventType, store.ErrGappedStream:
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			default:
				if err != nil {
					logger.Println(err)
					return err
				}
			}

			return c.JSON(res)
		case "get", "g":
			var req api.GetRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.Get(&req)
			switch err {
			default:
				if err != nil {
					logger.Println(err)
					return err
				}
			}

			return c.JSON(res)
		case "getall", "ga":
			var req api.GetAllRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.GetAll(&req)
			switch err {
			default:
				if err != nil {
					logger.Println(err)
					return err
				}
			}

			return c.JSON(res)
		case "eventcount", "ec":
			var req api.EventCountRequest
			// if err := json.Unmarshal(cmd.Payload, &req); err != nil {
			// 	return fiber.ErrUnprocessableEntity
			// }

			res, err := eventstore.EventCount(&req)
			switch err {
			default:
				if err != nil {
					logger.Println(err)
					return err
				}
			}

			return c.JSON(res)
		case "eventcountestimate", "ece":
			var req api.EventCountEstimateRequest
			// if err := json.Unmarshal(cmd.Payload, &req); err != nil {
			// 	return fiber.ErrUnprocessableEntity
			// }

			res, err := eventstore.EventCountEstimate(&req)
			switch err {
			default:
				if err != nil {
					logger.Println(err)
					return err
				}
			}

			return c.JSON(res)
		case "streamcount", "sc":
			var req api.StreamCountRequest
			// if err := json.Unmarshal(cmd.Payload, &req); err != nil {
			// 	return fiber.ErrUnprocessableEntity
			// }

			res, err := eventstore.StreamCount(&req)
			switch err {
			default:
				if err != nil {
					logger.Println(err)
					return err
				}
			}

			return c.JSON(res)
		case "streamcountestimate", "sce":
			var req api.StreamCountEstimateRequest
			// if err := json.Unmarshal(cmd.Payload, &req); err != nil {
			// 	return fiber.ErrUnprocessableEntity
			// }

			res, err := eventstore.StreamCountEstimate(&req)
			switch err {
			default:
				if err != nil {
					logger.Println(err)
					return err
				}
			}

			return c.JSON(res)
		case "liststreams", "ls":
			var req api.ListStreamsRequest
			if err := json.Unmarshal(cmd.Payload, &req); err != nil {
				return fiber.ErrUnprocessableEntity
			}

			res, err := eventstore.ListStreams(&req)
			switch err {
			default:
				if err != nil {
					logger.Println(err)
					return err
				}
			}

			return c.JSON(res)
		case "size", "s":
			var req api.SizeRequest
			// if err := json.Unmarshal(cmd.Payload, &req); err != nil {
			// 	return fiber.ErrUnprocessableEntity
			// }

			res, err := eventstore.Size(&req)
			switch err {
			default:
				if err != nil {
					logger.Println(err)
					return err
				}
			}

			return c.JSON(res)
		case "uptime", "up":
			// var req api.UptimeRequest
			// if err := json.Unmarshal(cmd.Payload, &req); err != nil {
			// 	return fiber.ErrUnprocessableEntity
			// }

			uptime := time.Since(start)

			return c.JSON(api.UptimeResponse{
				Uptime:      uptime.Milliseconds(),
				UptimeHuman: uptime.String(),
			})
		}

		return fiber.NewError(fiber.StatusBadRequest, "Unknown command")
	}
}

func RunHTTPServer(eventstore store.EventStore, logger *logrus.Logger) {
	httpPort := env.GetEnv("HTTP_PORT", "16543")
	tlsEnabled := env.GetEnv("TLS_ENABLED", "false") == "true"
	certFile := env.GetEnv("TLS_CERT_FILE", "certs/crt.pem")
	keyFile := env.GetEnv("TLS_KEY_FILE", "certs/key.pem")

	httpServer := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	httpServer.Post("/api/v1", HTTPHandler(eventstore, logger))

	go func() {
		logger.Printf("HTTP server listening on %s", httpPort)

		if tlsEnabled {
			if err := httpServer.ListenTLS(":"+httpPort, certFile, keyFile); err != nil {
				logger.Fatal(err)
			}
		} else {
			if err := httpServer.Listen(":" + httpPort); err != nil {
				logger.Fatal(err)
			}
		}
	}()
}

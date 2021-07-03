package transport

import (
	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/constants"
	"github.com/eventflowdb/eventflowdb/env"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func version() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(api.VersionResponse{
			Version: constants.Version,
		})
	}
}

func stream() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func all() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func add() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func streamCount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func event() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func eventCount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func uptime() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func streams() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func RunRestServer(eventstore store.EventStore, logger *logrus.Logger) *fiber.App {
	httpPort := env.GetEnv("HTTP_PORT", "16543")
	tlsEnabled := env.GetEnv("TLS_ENABLED", "false") == "true"
	certFile := env.GetEnv("TLS_CERT_FILE", "certs/crt.pem")
	keyFile := env.GetEnv("TLS_KEY_FILE", "certs/key.pem")

	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	v1 := server.Group("/api/v1")
	v1.Get("/version", version())
	v1.Get("/stream/all", all())
	v1.Get("/stream/count", streamCount())
	v1.Get("/stream/:id", stream())
	v1.Post("/stream/:id", add())
	v1.Get("/event/count", eventCount())
	v1.Get("/event/:id", event())
	v1.Get("/uptime", uptime())
	v1.Get("/streams", streams())

	go func() {
		logger.Printf("REST server listening on %s", httpPort)

		if tlsEnabled {
			if err := server.ListenTLS(":"+httpPort, certFile, keyFile); err != nil {
				logger.Fatal(err)
			}
		} else {
			if err := server.Listen(":" + httpPort); err != nil {
				logger.Fatal(err)
			}
		}
	}()

	return server
}

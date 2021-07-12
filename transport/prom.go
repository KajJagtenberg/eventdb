package transport

import (
	"github.com/eventflowdb/eventflowdb/env"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func RunPromServer(log *logrus.Logger) *fiber.App {
	promPort := env.GetEnv("PROM_PORT", "26543")

	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	server.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	go func() {
		log.Printf("Prometheus server listening on %s", promPort)

		if err := server.Listen(":" + promPort); err != nil {
			log.Fatal(err)
		}
	}()

	return server
}

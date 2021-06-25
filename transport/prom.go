package transport

import (
	"github.com/eventflowdb/eventflowdb/env"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func RunPromServer(logger *logrus.Logger) *fiber.App {
	promPort := env.GetEnv("PROM_PORT", "17654")
	tlsEnabled := env.GetEnv("TLS_ENABLED", "false") == "true"
	certFile := env.GetEnv("TLS_CERT_FILE", "certs/crt.pem")
	keyFile := env.GetEnv("TLS_KEY_FILE", "certs/key.pem")

	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	server.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	go func() {
		logger.Printf("Prometheus server listening on %s", promPort)

		if tlsEnabled {
			if err := server.ListenTLS(":"+promPort, certFile, keyFile); err != nil {
				logger.Fatal(err)
			}
		} else {
			if err := server.Listen(":" + promPort); err != nil {
				logger.Fatal(err)
			}
		}
	}()

	return server
}

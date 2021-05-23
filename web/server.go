package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kajjagtenberg/go-commando"
)

// var (
// 	//go:embed frontend/public/*
// 	frontend embed.FS
// )

type Options struct {
	Dispatcher *commando.CommandDispatcher
	Password   string
}

func CreateWebServer(options Options) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New())

	app.Post("/api/:cmd", func(c *fiber.Ctx) error {
		if len(options.Password) == 0 {
			return c.Next()
		}

		if c.Get("Authorization") != options.Password {
			return fiber.ErrUnauthorized
		}

		return c.Next()
	}, func(c *fiber.Ctx) error {
		result, err := options.Dispatcher.Handle(commando.Command{
			Name: c.Params("cmd"),
			Args: c.Body(),
		})
		if err != nil {
			return err
		}

		return c.JSON(result)
	})

	// f, err := fs.Sub(frontend, "frontend/public")
	// if err != nil {
	// 	return nil, err
	// }

	// app.Use(adaptor.HTTPHandler(http.FileServer(http.FS(f))))

	return app, nil
}

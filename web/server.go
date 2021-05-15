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

func CreateWebServer(dispatcher *commando.CommandDispatcher) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New())

	app.Post("/api", func(c *fiber.Ctx) error {
		var cmd commando.Command
		if err := c.BodyParser(&cmd); err != nil {
			return err
		}

		result, err := dispatcher.Handle(cmd)
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

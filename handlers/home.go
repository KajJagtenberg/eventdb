package handlers

import (
	"eventdb/constants"

	"github.com/gofiber/fiber/v2"
)

func Home() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{Name: constants.Name, Version: constants.Version})
	}
}

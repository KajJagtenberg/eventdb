package handlers

import (
	"eventdb/store"

	"github.com/gofiber/fiber/v2"
)

func GetEventCount(eventstore *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		count, err := eventstore.GetEventCount()
		if err != nil {
			return err
		}

		return c.JSON(struct {
			Count int `json:"count"`
		}{count})
	}
}

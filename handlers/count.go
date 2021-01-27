package handlers

import (
	"eventdb/store"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetEventCount(eventstore *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		count, err := eventstore.GetEventCount()
		if err != nil {
			log.Println(err)
			return fiber.ErrInternalServerError
		}

		return c.JSON(struct {
			Count int `json:"count"`
		}{count})
	}
}

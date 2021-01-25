package handlers

import (
	"eventdb/store"

	"github.com/gofiber/fiber/v2"
)

func Backup(eventstore *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Attachment("eventdb.bak")

		return eventstore.Backup(c.Response().BodyWriter())
	}
}

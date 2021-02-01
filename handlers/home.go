package handlers

import (
	"eventflowdb/constants"
	"eventflowdb/store"
	"eventflowdb/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Home(eventstore *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		size := eventstore.GetDBSize()

		eventCount, err := eventstore.GetEventCount()
		if err != nil {
			log.Println(err)

			return c.Status(fiber.StatusInternalServerError).JSON(Message{
				Message: "Internal server error",
			})
		}

		streamCount, err := eventstore.GetStreamCount()
		if err != nil {
			log.Println(err)

			return c.Status(fiber.StatusInternalServerError).JSON(Message{
				Message: "Internal server error",
			})
		}

		return c.JSON(struct {
			Name        string `json:"name"`
			Version     string `json:"version"`
			Size        int64  `json:"size"`
			HumanSize   string `json:"human_size"`
			EventCount  int    `json:"event_count"`
			StreamCount int    `json:"stream_count"`
		}{
			Name:        constants.Name,
			Version:     constants.Version,
			Size:        size,
			HumanSize:   util.ByteCountSI(size),
			EventCount:  eventCount,
			StreamCount: streamCount,
		})
	}
}

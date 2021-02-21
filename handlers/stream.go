package handlers

import (
	"log"
	"strconv"

	"eventflowdb/store"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

func LoadFromStream(eventstore *store.EventStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		streamParam := c.Params("stream")

		if len(streamParam) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Stream cannot be empty",
			})
		}

		stream, err := uuid.Parse(streamParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Stream must be an UUID v4",
			})
		}

		versionQuery := c.Query("version")
		limitQuery := c.Query("limit")

		version, _ := strconv.Atoi(versionQuery)
		limit, _ := strconv.Atoi(limitQuery)

		if version < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Version cannot be negative",
			})
		}

		if limit < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Limit cannot be negative",
			})
		}

		events, total, err := eventstore.LoadFromStream(stream, version, limit)
		if err != nil {
			log.Println(err)

			return c.Status(fiber.StatusInternalServerError).JSON(Message{
				Message: "Internal server error",
			})
		}

		return c.JSON(struct {
			Events  []store.RecordedEvent `json:"events"`
			Total   int                   `json:"total"`
			Version int                   `json:"version"`
			Limit   int                   `json:"limit"`
		}{
			events, total, version, limit,
		})
	}
}

func AppendToStream(eventstore *store.EventStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		streamParam := c.Params("stream")
		versionParam := c.Params("version")

		if len(streamParam) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Stream cannot be empty",
			})
		}

		stream, err := uuid.Parse(streamParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Stream must be an UUID v4",
			})
		}

		version, _ := strconv.Atoi(versionParam)

		if version < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Version cannot be negative",
			})
		}

		var body []struct {
			Type     string `json:"type" validate:"required"`
			Data     []byte `json:"data" validate:"required"`
			Metadata []byte `json:"metadata"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Unable to decode request body",
			})
		}

		if len(body) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Request contains no events",
			})
		}

		var events []store.EventData

		validator := validator.New()

		for _, event := range body {
			if err := validator.Struct(event); err != nil { // TODO: Add better error messages
				return c.Status(fiber.StatusBadRequest).JSON(Message{
					Message: err.Error(),
				})
			}

			events = append(events, store.EventData{
				Type:     event.Type,
				Data:     event.Data,
				Metadata: event.Metadata,
			})
		}

		if _, err := eventstore.AppendToStream(stream, version, events); err != nil {
			if err == store.ErrConcurrentStreamModification {
				return c.Status(fiber.StatusBadRequest).JSON(Message{
					Message: err.Error(),
				})
			} else {
				log.Println(err)

				return c.Status(fiber.StatusInternalServerError).JSON(Message{
					Message: "Internal server error",
				})
			}
		}

		return c.JSON(Message{
			Message: "Events added",
		})
	}
}

func Subscribe(eventstore *store.EventStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		offset := ulid.ULID{}

		events, err := eventstore.LoadFromAll(offset, 0)
		if err != nil {
			log.Println(err)
			return fiber.ErrInternalServerError
		}

		return c.JSON(events)
	}
}

func GetEventByID(eventstore *store.EventStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("Getting the event from the store")

		idParam := c.Params("id")

		if len(idParam) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Event ID cannot be empty",
			})
		}

		id, err := ulid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: err.Error(),
			})
		}

		event, err := eventstore.GetEventByID(id)
		if err != nil {
			log.Println(err)

			return c.Status(fiber.StatusInternalServerError).JSON(Message{
				Message: "Internal server error",
			})
		}

		return c.JSON(event)
	}
}

func GetStreams(eventstore *store.EventStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		offsetQuery := c.Query("offset")
		limitQuery := c.Query("limit")

		offset, _ := strconv.Atoi(offsetQuery)
		limit, _ := strconv.Atoi(limitQuery)

		if offset < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Offset cannot be negative",
			})
		}

		if limit < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(Message{
				Message: "Limit cannot be negative",
			})
		}

		streams, total, err := eventstore.GetStreams(offset, limit)

		if err != nil {
			log.Println(err)

			return c.Status(fiber.StatusInternalServerError).JSON(Message{
				Message: "Internal server error",
			})
		}

		return c.JSON(struct {
			Streams []uuid.UUID `json:"streams"`
			Total   int         `json:"total"`
			Offset  int         `json:"offset"`
			Limit   int         `json:"limit"`
		}{
			streams, total, offset, limit,
		})
	}
}

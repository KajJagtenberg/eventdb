package handlers

import (
	"errors"
	"log"
	"strconv"

	"eventdb/store"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

func LoadFromStream(eventstore *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		streamParam := c.Params("stream")

		if len(streamParam) == 0 {
			return errors.New("Stream cannot be empty")
		}

		stream, err := uuid.Parse(streamParam)
		if err != nil {
			return errors.New("Stream must be an UUID v4")
		}

		versionQuery := c.Query("version")
		limitQuery := c.Query("limit")

		version, _ := strconv.Atoi(versionQuery)
		limit, _ := strconv.Atoi(limitQuery)

		if version < 0 {
			return errors.New("Version cannot be negative")
		}

		if limit < 0 {
			return errors.New("Limit cannot be negative")
		}

		events, total, err := eventstore.LoadFromStream(stream, version, limit)
		if err != nil {
			return err
		}

		if len(events) == 0 {
			return errors.New("Not Found")
		}

		return c.JSON(struct {
			Streams []store.Event `json:"events"`
			Total   int           `json:"total"`
			Version int           `json:"version"`
			Limit   int           `json:"limit"`
		}{
			events, total, version, limit,
		})
	}
}

func AppendToStream(eventstore *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		streamParam := c.Params("stream")
		versionParam := c.Params("version")

		if len(streamParam) == 0 {
			return errors.New("Stream cannot be empty")
		}

		stream, err := uuid.Parse(streamParam)
		if err != nil {
			return errors.New("Stream must be an UUID v4")
		}

		version, _ := strconv.Atoi(versionParam)

		if version < 0 {
			return errors.New("Version cannot be negative")
		}

		var events []store.AppendEvent

		if err := c.BodyParser(&events); err != nil {
			return err
		}

		if len(events) == 0 {
			return errors.New("Empty events")
		}

		validate := validator.New()

		for _, event := range events {
			if err := validate.Struct(event); err != nil {
				return err
			}
		}

		if err := eventstore.AppendToStream(stream, version, events); err != nil {
			return err
		}

		return c.SendString("Events added")
	}
}

func Subscribe(eventstore *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		offset := ulid.ULID{}

		events, err := eventstore.Subscribe(offset, 0)
		if err != nil {
			return err
		}

		return c.JSON(events)
	}
}

func GetEventByID(eventstore *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("Getting the event from the store")

		idParam := c.Params("id")

		if len(idParam) == 0 {
			return errors.New("Event ID cannot be empty")
		}

		id, err := ulid.Parse(idParam)
		if err != nil {
			return err
		}

		event, err := eventstore.GetEventByID(id)
		if err != nil {
			return err
		}

		return c.JSON(event)
	}
}

func GetStreams(eventstore *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		offsetQuery := c.Query("offset")
		limitQuery := c.Query("limit")

		offset, _ := strconv.Atoi(offsetQuery)
		limit, _ := strconv.Atoi(limitQuery)

		if offset < 0 {
			return errors.New("Offset cannot be negative")
		}

		if limit < 0 {
			return errors.New("Limit cannot be negative")
		}

		streams, total, err := eventstore.GetStreams(offset, limit)

		if err != nil {
			log.Println(err)
			return err
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

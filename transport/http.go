package transport

import (
	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/constants"
	"github.com/gofiber/fiber/v2"
)

func SetupHandler(r fiber.Router) {
	r.Post("/version", VersionHandler)
}

func VersionHandler(c *fiber.Ctx) error {
	var req api.VersionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	return c.JSON(api.VersionResponse{
		Version: constants.Version,
	})
}

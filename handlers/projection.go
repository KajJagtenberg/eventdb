package handlers

import "github.com/gofiber/fiber/v2"

// POST /api/v1/projections - Create a new projection

func CreateProjection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

// GET /api/v1/projections - Returns list of projections

func GetProjections() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

// GET /api/v1/projections/{name} - Returns info about a single projection

func GetProjection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

// GET /api/v1/projections/{name}/{collection} - Returns a list of documents from a collection

func GetCollection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

// GET /api/v1/projections/{name}/{colleciton}/{id} - Returns a single document from a collection

func GetDocument() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

// PUT /api/v1/projections/{name} - Updates a projection, restarting it if neccesary

func UpdateProjection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

// DELETE /api/v1/projections/{name} - Deletes a projection with all of its collections and documents

func DeleteProjection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

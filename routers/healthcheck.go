package routers

import (
	"github.com/gofiber/fiber/v2"
)

func healthCheck(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is healthy",
		})
	})
}

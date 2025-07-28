package routers

import "github.com/gofiber/fiber/v2"

func StartApp() {
	app := fiber.New()

	healthCheck(app)

	app.Listen(":3000")
}

package routers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func StartApp() {
	app := fiber.New()

	healthCheck(app)

	// Get port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen("[::]:" + port))
}

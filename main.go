package main

import (
	"log"
	"os"

	"github.com/Logan9312/Prize-Bot-V2/routers"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	routers.HealthCheck(app)

	// Get port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen("[::]:" + port))
}

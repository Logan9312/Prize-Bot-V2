package main

import (
	"fmt"
	"log"

	"github.com/Logan9312/Prize-Bot-V2/routers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Bot is running!")

	app := fiber.New()

	routers.HealthCheck(app)

	log.Fatal(app.Listen(":3000"))
}

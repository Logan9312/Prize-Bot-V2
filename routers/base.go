package routers

import (
	"os"

	"github.com/Logan9312/Prize-Bot-V2/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func StartApp() {
	app := fiber.New()

	healthCheck(app)

	// Get port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	logger.Logger.Info("Starting web server", zap.String("port", port))
	if err := app.Listen("[::]:" + port); err != nil {
		logger.Logger.Fatal("Failed to start web server", zap.Error(err))
	}
}

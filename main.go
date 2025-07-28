package main

import (
	"log"

	"github.com/Logan9312/Prize-Bot-V2/bot"
	"github.com/Logan9312/Prize-Bot-V2/logger"
	"github.com/Logan9312/Prize-Bot-V2/routers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if err := logger.Init(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Close()

	bot.Start()
	routers.StartApp()
}

package main

import (
	"fmt"
	"os"

	"github.com/Logan9312/Prize-Bot-V2/connect"
	"github.com/Logan9312/Prize-Bot-V2/database"
	"github.com/Logan9312/Prize-Bot-V2/routers"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v72"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	fmt.Println(os.Environ())

	stripe.Key = os.Getenv("STRIPE_TOKEN")

	//Connects database
	database.DatabaseConnect(os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("ENVIRONMENT"))

	//Connects main bot
	mainSession, err := connect.BotConnect(os.Getenv("DISCORD_TOKEN"), os.Getenv("ENVIRONMENT"))
	if err != nil {
		fmt.Println(err)
	}

	devData := database.DevSetup{
		BotID: mainSession.State.User.ID,
	}
	result := database.DB.First(&devData)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	WhiteLabels := []map[string]any{}
	result = database.DB.Model([]database.WhiteLabels{}).Find(&WhiteLabels)
	if result.Error != nil {
		fmt.Println("Error fetching whitelabels:", result.Error)
	}

	for _, v := range WhiteLabels {
		s, err := connect.BotConnect(v["bot_token"].(string), os.Getenv("ENVIRONMENT"))
		if err != nil {
			fmt.Println("Error connecting bot: %w", err)
		}
		if s.State.User.ID == "995022149226082324" {
			err = s.UpdateGameStatus(0, "Bot Version "+devData.Version)
			if err != nil {
				fmt.Println("Error setting status", err)
			}
		}
	}

	//go commands.SetRoles(mainSession)

	fmt.Println("Bot is running!")

	routers.HealthCheck()

}

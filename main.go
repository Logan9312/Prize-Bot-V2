package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Logan9312/Prize-Bot-V2/connect"
	"github.com/Logan9312/Prize-Bot-V2/database"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v72"
)

// Environment struct
type Environment struct {
	Environment  string `env:"ENVIRONMENT,required"`
	DiscordToken string `env:"DISCORD_TOKEN,required"`
	Host         string `env:"DB_HOST"`
	Password     string `env:"DB_PASSWORD"`
	StripeToken  string `env:"STRIPE_TOKEN"`
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	value := os.Getenv("DISCORD_TOKEN")
	fmt.Println("Value:", value)

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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

}

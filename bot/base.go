package bot

import (
	"os"

	"github.com/Logan9312/Prize-Bot-V2/logger"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func Start() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		logger.Logger.Fatal("Failed to create Discord session", zap.Error(err))
	}

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Logger.Info("Bot logged in", 
			zap.String("username", r.User.Username), 
			zap.String("discriminator", r.User.Discriminator))
	})

	err = dg.Open()
	if err != nil {
		logger.Logger.Fatal("Failed to open Discord connection", zap.Error(err))
	}
}

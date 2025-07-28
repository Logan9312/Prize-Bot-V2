package bot

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func Start() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", r.User.Username, r.User.Discriminator)
	})

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
}

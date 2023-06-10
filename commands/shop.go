package commands

import (
	. "github.com/Logan9312/Prize-Bot-V2/models"
	"github.com/bwmarrin/discordgo"
)

var ShopCommand = discordgo.ApplicationCommand{
	Name:        "shop",
	Description: "Buy, Sell and Trade items with others!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "sell",
			Description: "Select an item to sell",
			Options: []*discordgo.ApplicationCommandOption{
				CommandOptionItem,
			},
		},
	},
}

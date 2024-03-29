package commands

import (
	
	"github.com/bwmarrin/discordgo"
	u "github.com/Logan9312/Prize-Bot-V2/utils"
)

var CurrencyCommand = discordgo.ApplicationCommand{
	Name:        "currency",
	Description: "Shows the current currency of the bot",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Adds to existing currency amount.",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionTarget),
				Require(CommandOptionAmount),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Sets currency to a specific value.",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionTarget),
				Require(CommandOptionAmount),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "Lists the currency of all saved users.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "display",
			Description: "Display the currency of an user.",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionUser),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "manage",
			Description: "Manage the currencies in your server.",
		},
	},
	DefaultMemberPermissions: u.Ptr(int64(discordgo.PermissionManageServer)),
	DMPermission:             new(bool),
}

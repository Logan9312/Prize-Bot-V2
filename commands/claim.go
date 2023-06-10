package commands

import (
	"github.com/Logan9312/Prize-Bot-V2/helpers"
	"github.com/bwmarrin/discordgo"
)

var ClaimCommand = discordgo.ApplicationCommand{
	Name:        "claim",
	Description: "Manage your claims!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Create a claimable prize.",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionTarget),
				Require(CommandOptionItem),
				Require(CommandOptionLogChannel),
				CommandOptionDescription,
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "inventory",
			Description: "Displays a user's unclaimed prizes",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionUser),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "refresh",
			Description: "Resends all claim messages. Will not delete old ones",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionChannel),
			},
		},
	},
	DefaultMemberPermissions: helpers.Ptr(int64(discordgo.PermissionManageServer)),
}

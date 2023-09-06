package commands

import (
	u "github.com/Logan9312/Prize-Bot-V2/utils"
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
	DefaultMemberPermissions: u.Ptr(int64(discordgo.PermissionManageServer)),
}

func ClaimOutput() {
	//TODO Move this to claim output
	/* 	if auction.TargetPrice != nil && *auction.TargetPrice > auction.Bid {

		auctionMap["target_message"] = fmt.Sprintf("The host had set a target price of %s that has not been reached.", PriceFormat(auctionMap["target_price"].(float64), guildID, auctionMap["currency"]))
	} */

	//TODO Add in a message about this when the auction ends
	//TODO Move this to the claiming process.
	/*
		if auctionMap["use_currency"] != nil && auctionMap["use_currency"].(bool) && auctionMap["winner"] != nil {
			err = CurrencySubtractUser(guildID, auctionMap["winner"].(string), auctionMap["bid"].(float64))
			if err != nil {
				return err
			}
			err = CurrencyAddUser(guildID, auctionMap["host"].(string), auctionMap["bid"].(float64))
			if err != nil {
				return err
			}

		} */

	//TODO Potentially move this to claiming
	/* 	if auction.Buyout != nil {
		if auction.Bid >= *auction.Buyout {
			auctionMap["buyout_message"] = fmt.Sprintf("%s\n\u200bBUYOUT!", PriceFormat(auctionMap["buyout"].(float64), auctionMap["guild_id"].(string), auctionMap["currency"]))
		}
	} */
}

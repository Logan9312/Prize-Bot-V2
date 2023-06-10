package connect

import (
	"fmt"
	"strings"
	"time"

	c "github.com/Logan9312/Prize-Bot-V2/commands"
	r "github.com/Logan9312/Prize-Bot-V2/responses"
	"github.com/bwmarrin/discordgo"
)

// Move these to commands package
var commandMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error{
	"auction": c.Auction,
	/*"bid":            c.AuctionBid,
	"profile":        c.Profile,
	"giveaway":       c.Giveaway,
	"shop":           c.Shop,
	"claim":          c.Claim,
	"privacy_policy": c.Privacy,
	"dev":            c.Dev,
	"premium":        c.Premium,
	"settings":       c.Settings,
	"quest":          c.Quest,
	"currency":       c.Currency,
	"stats":          c.Stats,
	"help":           Help,
	"whitelabel":     Whitelabel,*/
}

var buttonMap = map[string]func(*discordgo.Session, *discordgo.InteractionCreate) error{
	/*	"endauction":             c.AuctionEndButton,
		"claim_prize":            c.ClaimPrizeButton,
		"clearauction":           c.ClearAuctionButton,
		"delete_auction_queue":   c.DeleteAuctionQueue,
		"delete_auction_channel": c.DeleteAuctionChannel,
		"reroll_giveaway":        c.RerollGiveawayButton,
		"clear_auction_setup":    c.AuctionSetupClearButton,
		"clear_giveaway_setup":   c.GiveawaySetupClearButton,
		"clear_claim_setup":      c.ClaimSetupClearButton,
		"clear_currency_setup":   c.CurrencySetupClearButton,
		"claim_cancel":           c.CancelButton,
		"claim_complete":         c.CompleteButton,
		"reopen_ticket":          c.ReopenTicket,
		"additem":                c.AddItem,
		"bid_history":            c.AuctionBidHistory,
		"questbutton":            c.QuestButton,*/
}

var autoCompleteMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error{
	/*	"auction":  c.AuctionAutoComplete,
		"giveaway": c.GiveawayAutoComplete,*/
}

var modalSubmitMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error{}

var guildMembersMap = map[string]func(s *discordgo.Session, g *discordgo.GuildMembersChunk) error{
	/*"claim_create": c.ClaimCreateRole,
	"$":            c.CurrencyRoleHandler,*/
}

func RegisterHandlers(s *discordgo.Session) {
	s.AddHandler(ReadyHandler)
	s.AddHandler(InteractionHandler)
	//s.AddHandler(MessageHandler)
	//s.AddHandler(GuildMemberChunkHandler)
	s.AddHandler(GuildCreateHandler)
}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		r.ErrorResponse(s, i, fmt.Errorf("commands cannot be run in a DM. Please contact support if you're not currently in a DM with the bot"))
		return
	}

	if f := InteractionRouter(s, i); f != nil {
		if err := f(s, i); err != nil {
			r.ErrorResponse(s, i, err)
		}
	} else {
		r.ErrorResponse(s, i, fmt.Errorf("response has not been set properly, please contact Logan to fix"))
	}

}

func InteractionRouter(s *discordgo.Session, i *discordgo.InteractionCreate) func(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		fmt.Println(i.ApplicationCommandData().Name, "is being run by:", i.Member.User.Username)
		return commandMap[i.ApplicationCommandData().Name]

	case discordgo.InteractionMessageComponent:
		fmt.Println(i.MessageComponentData().CustomID, "is being run by:", i.Member.User.Username)
		return buttonMap[strings.Split(i.MessageComponentData().CustomID, ":")[0]]

	case discordgo.InteractionApplicationCommandAutocomplete:
		return autoCompleteMap[i.ApplicationCommandData().Name]

	case discordgo.InteractionModalSubmit:
		fmt.Println(i.ModalSubmitData().CustomID, "is being run by:", i.Member.User.Username)
		if i.ModalSubmitData().CustomID == "whitelabel_token" {
			return WhitelabelTokenModal
		} else {
			return modalSubmitMap[i.ModalSubmitData().CustomID]
		}
	}

	return nil
}

func ReadyHandler(s *discordgo.Session, i *discordgo.Ready) {
	_, err := s.ChannelMessageSend("943175605858496602", "Bot has finished restarting")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Bot is ready")
}

func GuildCreateHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	channelID := "1005255087200948234"

	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Name",
			Value: g.Name,
		},
		{
			Name:  "Member Count",
			Value: fmt.Sprint(g.MemberCount),
		},
	}
	if g.VanityURLCode != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Vanity URL",
			Value: g.VanityURLCode,
		})
	}

	if !g.Unavailable {
		_, err := r.SuccessMessage(s, channelID, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "New Server Joined!",
					Description: "The bot has joined a new discord within the last 5 minutes.",
					Fields:      fields,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: g.IconURL(""),
					},
					Image: &discordgo.MessageEmbedImage{
						URL: g.BannerURL(""),
					},
				},
			},
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

/*func GuildMemberChunkHandler(s *discordgo.Session, g *discordgo.GuildMembersChunk) {
	fmt.Println(g.Nonce)
	if f, ok := guildMembersMap[strings.Split(g.Nonce, ":")[0]]; ok {
		err := f(s, g)
		if err != nil {
			fmt.Println(err)
			_, err = FollowUpErrorResponse(s, ReadChunkData(strings.Split(g.Nonce, ":")[1])["interaction"].(*discordgo.InteractionCreate), err.Error())
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("GuildMemberChunkHandler has not been set properly.")
	}
}*/

func DeleteCommandMessage(s *discordgo.Session, channelID, messageID string) error {
	time.Sleep(10 * time.Second)
	return s.ChannelMessageDelete(channelID, messageID)
}

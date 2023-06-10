package commands

import (
	"fmt"
	"strings"

	"github.com/Logan9312/Prize-Bot-V2/database"
	"github.com/Logan9312/Prize-Bot-V2/events"
	h "github.com/Logan9312/Prize-Bot-V2/helpers"
	r "github.com/Logan9312/Prize-Bot-V2/responses"
	"github.com/bwmarrin/discordgo"
)

var AuctionCommand = discordgo.ApplicationCommand{
	Name:        "auction",
	Description: "Put an item up for auction!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "schedule",
			Description: "Schedule an auction",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionTimeUntil),
				Require(CommandOptionItem),
				Require(CommandOptionBid),
				Require(CommandOptionDuration),
				CommandOptionDescription,
				CommandOptionCurrency,
				CommandOptionUseCurrency,
				CommandOptionIntegerOnly,
				CommandOptionIncrementMax,
				CommandOptionIncrementMin,
				CommandOptionBuyout,
				CommandOptionTargetPrice,
				CommandOptionCategory,
				CommandOptionImage,
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Create an Auction",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionItem),
				Require(CommandOptionBid),
				Require(CommandOptionDuration),
				CommandOptionDescription,
				CommandOptionCurrency,
				CommandOptionUseCurrency,
				CommandOptionIntegerOnly,
				CommandOptionIncrementMax,
				CommandOptionIncrementMin,
				CommandOptionBuyout,
				CommandOptionTargetPrice,
				CommandOptionCategory,
				CommandOptionImage,
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "queue",
			Description: "Display the current auction Queue",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "edit",
			Description: "Edit any auction details",
			Options: []*discordgo.ApplicationCommandOption{
				//CommandOptionQueueNumber
				CommandOptionItem,
				CommandOptionHost,
				CommandOptionBid,
				CommandOptionWinner,
				CommandOptionExtend,
				CommandOptionDescription,
				CommandOptionCurrency,
				CommandOptionUseCurrency,
				CommandOptionCurrencySide,
				CommandOptionIntegerOnly,
				CommandOptionIncrementMax,
				CommandOptionIncrementMin,
				CommandOptionBuyout,
				CommandOptionTargetPrice,
				CommandOptionImage,
			},
		},
	},
}

var BidCommand = discordgo.ApplicationCommand{
	Name:        "bid",
	Description: "Bid on an Auction",
	Options: []*discordgo.ApplicationCommandOption{
		Require(CommandOptionBid),
	},
}

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Options[0].Name {
	case "schedule":
		//return AuctionSchedule(s, i)
	case "create":
		return AuctionCreate(s, i)
	case "queue":
		//return AuctionQueue(s, i)
	case "edit":
		//return AuctionEdit(s, i)
	}
	return fmt.Errorf("unknown Auction command, please contact support")
}

func SaveAuction(auction database.Auction) error {
	//TODO wtf is gorm save??
	return database.DB.Save(&auction).Error
}

func AuctionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := h.ParseCommand(i)
	errors := []string{}

	//Splits the auction names
	auctions := strings.Split(options["item"].(string), ";")

	if len(auctions) > 5 && !CheckPremiumGuild(i.GuildID) {
		r.PremiumError(s, i, "Free users can only start 5 auctions in bulk. Upgrade to premium to start up to 100 in bulk.")
	}

	if len(auctions) > 100 {
		return fmt.Errorf("you can only start 100 auctions in bulk at once. You attempted to start: %d", len(auctions))
	}

	auctionData, err := SetAuctionData(s, i, options)
	if err != nil {
		return fmt.Errorf("error setting auction data: %w", err)
	}

	//TODO Optimize selecting multiple auctions
	for _, item := range auctions {
		//Copies all of the options to a fresh auctionMap
		auctionData.Event.Item = item

		err := SaveAuction(auctionData)
		if err != nil {
			errors = append(errors, err.Error())
		}

		channelID, err := AuctionStart(s, auctionData)
		if err != nil {
			errors = append(errors, err.Error())
		}

		_, err = r.FollowUpSuccessResponse(s, i, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "**Auction Starting**",
					Description: fmt.Sprintf("Auction has successfully been started in <#%s>!", channelID),
				},
			},
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("one or more auctions failed to start:\n%s", strings.Join(errors, "\n"))
	}
	return nil
}

func SetAuctionData(s *discordgo.Session, i *discordgo.InteractionCreate, data map[string]any) (database.Auction, error) {

	event, err := events.NewEvent(s.State.User.ID, i, data)
	if err != nil {
		return database.Auction{}, r.Errorfmt(err)
	}

	auctionData := database.Auction{
		Event:    event,
		Bid:      data["bid"].(float64),
		Currency: &database.Currency{},
	}

	if data["currency"] != nil {
		//TODO Handle Currency
	}

	if data["increment_min"] != nil {
		auctionData.IncrementMin = h.Ptr(data["increment_min"].(float64))
	}

	if data["increment_max"] != nil {
		auctionData.IncrementMax = h.Ptr(data["increment_max"].(float64))
	}

	if data["target_price"] != nil {
		auctionData.TargetPrice = h.Ptr(data["target_price"].(float64))
	}

	if data["buyout"] != nil {
		auctionData.Buyout = h.Ptr(data["buyout"].(float64))
	}

	if data["integer_only"] != nil {
		auctionData.IntegerOnly = data["integer_only"].(bool)
	}

	if data["channel_lock"] != nil {
		auctionData.ChannelLock = h.Ptr(data["channel_lock"].(bool))
	}

	return auctionData, nil

}

func AuctionHostCheck(auctionSetup map[string]any, member *discordgo.Member) bool {
	if auctionSetup["host_role"] == nil {
		return true
	}
	for _, v := range member.Roles {
		if v == auctionSetup["host_role"].(string) {
			return true
		}
	}

	return member.Permissions&discordgo.PermissionManageServer == 8
}

func AuctionStart(s *discordgo.Session, data database.Auction) (string, error) {

	auctionSettings := &database.AuctionSetup{}
	currencySettings := &database.CurrencySetup{}

	result := database.DB.First(&auctionSettings, data.Event.GuildID)
	if result.Error != nil {
		fmt.Println("Error getting auction settings: " + result.Error.Error())
	}

	result = database.DB.First(&currencySettings, data.Event.GuildID)
	if result.Error != nil {
		fmt.Println("Error getting currency settings: " + result.Error.Error())
	}

	//TODO Mark queued auction as started in database
	auctionMessage := AuctionMessageFormat(data)

	//TODO Make prefix work for editing and fix prefix working on create channel.

	//TODO Rebuild Channel Lock
	if !*data.ChannelLock {
		channelCreateData := discordgo.GuildChannelCreateData{
			Name: fmt.Sprintf("%s%s", *auctionSettings.ChannelPrefix, data.Event.Item),
			Type: discordgo.ChannelTypeGuildText,
		}
		if auctionSettings.Category != nil {
			channelCreateData.ParentID = *auctionSettings.Category
		}

		channel, err := s.GuildChannelCreateComplex(data.Event.GuildID, channelCreateData)
		if err != nil {
			return "", err
		}
		data.Event.ChannelID = &channel.ID
	}

	message, err := r.SuccessMessage(s, *data.Event.ChannelID, &auctionMessage)
	if err != nil {
		return *data.Event.ChannelID, err
	}

	data.Event.MessageID = &message.ID

	err = SaveAuction(data)
	if err != nil {
		return *data.Event.ChannelID, fmt.Errorf("error saving auction to database, auction will not work: %w", result.Error)
	}

	return *data.Event.ChannelID, nil
}

func AuctionMessageFormat(data database.Auction) discordgo.MessageSend {
	message := events.MessageFormat(data.Event)

	if data.IncrementMin != nil {
		message.Embeds[0].Description += fmt.Sprintf("**Minimum Bid:** + %s above previous.\n", h.PriceFormat(*data.IncrementMin, data.Event.GuildID, data.Currency))
	}

	if data.IncrementMax != nil {
		message.Embeds[0].Description += fmt.Sprintf("**Maximum Bid:** + %s above previous.\n", h.PriceFormat(*data.IncrementMax, data.Event.GuildID, data.Currency))
	}

	if data.TargetPrice != nil {
		message.Embeds[0].Description += "**Target Price:** The host has set a hidden target price for this auction.\n"
	}

	if data.IntegerOnly {
		message.Embeds[0].Description += fmt.Sprintf("**Integer Only:** %t.\n", data.IntegerOnly)
	}

	if data.SnipeExtension != nil && data.SnipeRange != nil {
		message.Embeds[0].Description += fmt.Sprintf("**Anti Snipe:** If a bid is placed within the last %s, the auction will be extended by %s.\n", data.SnipeRange, data.SnipeExtension.String())
	}

	if data.Buyout != nil {
		message.Embeds[0].Description += fmt.Sprintf("**Buyout Price:** %s.\n", h.PriceFormat(*data.Buyout, data.Event.GuildID, data.Currency))
	}

	fieldName := "__**Starting Bid:**__"
	if data.WinnerID != nil {
		fieldName = "__**Current Highest Bid:**__"
	}

	message.Embeds[0].Fields = append(message.Embeds[0].Fields, &discordgo.MessageEmbedField{
		Name:  "__**How to Bid**__",
		Value: "Use the /bid command or type `/bid {value}` in chat\nEx: `/bid 550`\n",
	})

	message.Embeds[0].Fields = append(message.Embeds[0].Fields, &discordgo.MessageEmbedField{
		Name:   fieldName,
		Value:  h.PriceFormat(data.Bid, data.Event.GuildID, data.Currency),
		Inline: true,
	})

	if data.WinnerID != nil {
		message.Embeds[0].Fields = append(message.Embeds[0].Fields, &discordgo.MessageEmbedField{
			Name:   "__**Current Winner**__",
			Value:  fmt.Sprintf("<@%s>", *data.WinnerID),
			Inline: true,
		})
	}

	message.Components = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "End Auction",
					Style: 4,
					Emoji: discordgo.ComponentEmoji{
						Name: "🛑",
					},
					CustomID: "endauction",
				},
				discordgo.Button{
					Label:    "Clear Chat",
					Style:    3,
					CustomID: "clearauction",
					Emoji: discordgo.ComponentEmoji{
						Name: "restart",
						ID:   "835685528917114891",
					},
					Disabled: false,
				},
			},
		},
	}

	if data.BidHistory != nil {
		if len(*data.BidHistory) > 4095 {
			data.BidHistory = h.Ptr((*data.BidHistory)[len(*data.BidHistory)-4095:])
		}
		message.Embeds = append(message.Embeds, &discordgo.MessageEmbed{
			Title:       "**Bid History**",
			Description: *data.BidHistory,
			Color:       0x8073ff,
			Image: &discordgo.MessageEmbedImage{
				URL: "https://i.imgur.com/9wo7diC.png",
			},
		})
	}

	return message
}

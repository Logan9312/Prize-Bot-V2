package commands

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Logan9312/Prize-Bot-V2/database"
	"github.com/Logan9312/Prize-Bot-V2/events"
	u "github.com/Logan9312/Prize-Bot-V2/utils"
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

var BuyoutCommand = discordgo.ApplicationCommand{
	Name:        "buyout",
	Description: "Buyout an Auction",
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

func AuctionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := u.ParseSubCommand(i)
	errors := []string{}

	//Splits the auction names
	auctions := strings.Split(options["item"].(string), ";")

	if len(auctions) > 5 && !CheckPremiumGuild(i.GuildID) {
		u.PremiumError(s, i, "Free users can only start 5 auctions in bulk. Upgrade to premium to start up to 100 in bulk.")
	}

	if len(auctions) > 100 {
		return fmt.Errorf("you can only start 100 auctions in bulk at once. You attempted to start: %d", len(auctions))
	}

	auctionData, err := SetAuctionData(s, i, options)
	if err != nil {
		return fmt.Errorf("error setting auction data: %w", err)
	}

	u.SuccessResponse(s, i, &discordgo.InteractionResponseData{
		Content:    "",
		Components: []discordgo.MessageComponent{},
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "**Starting Auctions**",
				Description: "Please wait while the auctions are being created",
			},
		},
	})

	//TODO Optimize selecting multiple auctions
	for _, item := range auctions {
		//Copies all of the options to a fresh auctionMap
		auctionData.Event.Item = item
		err := database.SaveAuction(&auctionData)
		if err != nil {
			errors = append(errors, err.Error())
		}

		channelID, err := AuctionStart(s, auctionData)
		if err != nil {
			errors = append(errors, err.Error())
		}

		_, err = u.FollowUpSuccessResponse(s, i, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "**Auction Started**",
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
		return database.Auction{}, u.Errorfmt(err)
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
		auctionData.IncrementMin = u.Ptr(data["increment_min"].(float64))
	}

	if data["increment_max"] != nil {
		auctionData.IncrementMax = u.Ptr(data["increment_max"].(float64))
	}

	if data["target_price"] != nil {
		auctionData.TargetPrice = u.Ptr(data["target_price"].(float64))
	}

	if data["buyout"] != nil {
		auctionData.Buyout = u.Ptr(data["buyout"].(float64))
	}

	if data["integer_only"] != nil {
		auctionData.IntegerOnly = data["integer_only"].(bool)
	}

	//TODO Rework Channel Lock so it's actually useful

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

	currencySettings := &database.CurrencySetup{}

	auctionSettings, err := database.GetAuctionSettings(data.Event.GuildID)
	if err != nil {
		fmt.Println(fmt.Errorf("error getting auction settings: %w", err))
	}

	result := database.DB.First(&currencySettings, data.Event.GuildID)
	if result.Error != nil {
		fmt.Println("Error getting currency settings: " + result.Error.Error())
	}

	//TODO Mark queued auction as started in database
	auctionMessage := AuctionMessageFormat(data)

	//TODO Rebuild Channel Lock
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

	message, err := u.SuccessMessage(s, *data.Event.ChannelID, &auctionMessage)
	if err != nil {
		return *data.Event.ChannelID, err
	}

	data.Event.MessageID = &message.ID
	err = database.SaveAuction(&data)
	if err != nil {
		return *data.Event.ChannelID, fmt.Errorf("error saving auction to database, auction will not work: %w", result.Error)
	}

	fmt.Print()

	return *data.Event.ChannelID, nil
}

func AuctionMessageFormat(data database.Auction) discordgo.MessageSend {
	message := events.MessageFormat(data.Event)

	if data.IncrementMin != nil {
		message.Embeds[0].Fields[0].Value += fmt.Sprintf("**Minimum Bid:** + %s above previous.\n", u.PriceFormat(*data.IncrementMin, data.Event.GuildID, data.Currency))
	}

	if data.IncrementMax != nil {
		message.Embeds[0].Fields[0].Value += fmt.Sprintf("**Maximum Bid:** + %s above previous.\n", u.PriceFormat(*data.IncrementMax, data.Event.GuildID, data.Currency))
	}

	if data.TargetPrice != nil {
		message.Embeds[0].Fields[0].Value += "**Target Price:** The host has set a hidden target price for this auction.\n"
	}

	if data.IntegerOnly {
		message.Embeds[0].Fields[0].Value += fmt.Sprintf("**Integer Only:** %t.\n", data.IntegerOnly)
	}

	if data.SnipeExtension != nil && data.SnipeRange != nil {
		message.Embeds[0].Fields[0].Value += fmt.Sprintf("**Anti Snipe:** If a bid is placed within the last %s, the auction will be extended by %s.\n", data.SnipeRange, data.SnipeExtension.String())
	}

	if data.Buyout != nil {
		message.Embeds[0].Fields[0].Value += fmt.Sprintf("**Buyout Price:** %s.\n", u.PriceFormat(*data.Buyout, data.Event.GuildID, data.Currency))
	}

	fieldName := "__**Starting Bid:**__"
	if data.WinnerID != nil {
		fieldName = "__**Current Highest Bid:**__"
	}

	message.Embeds[0].Fields = append(message.Embeds[0].Fields, &discordgo.MessageEmbedField{
		Name:   fieldName,
		Value:  u.PriceFormat(data.Bid, data.Event.GuildID, data.Currency),
		Inline: true,
	})

	if data.WinnerID != nil {
		message.Embeds[0].Fields = append(message.Embeds[0].Fields, &discordgo.MessageEmbedField{
			Name:   "__**Current Winner**__",
			Value:  fmt.Sprintf("<@%s>", *data.WinnerID),
			Inline: true,
		})
	}

	message.Embeds[0].Fields = append(message.Embeds[0].Fields, &discordgo.MessageEmbedField{
		Name:  "__**How to Bid**__",
		Value: "Use the /bid command or type `/bid {value}` in chat\nEx: `/bid 550`\n",
	})

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
			data.BidHistory = u.Ptr((*data.BidHistory)[len(*data.BidHistory)-4095:])
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

func AuctionBid(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := u.ParseSlashCommand(i)
	bid := options["bid"].(float64)

	auction, err := database.GetAuctionData(i.ChannelID)
	if err != nil {
		return err
	}

	//FIXME This was just a quick fix to prevent both bots trying to place a bit. Rework if needed once I can save auctions with bot ID
	message, err := s.ChannelMessage(i.ChannelID, *auction.Event.MessageID)
	if err != nil {
		return err
	}
	if message.Author.ID != s.State.User.ID {
		return fmt.Errorf("You must use the bot that started the auction to place a bid.")
	}

	if auction.SnipeRange != nil && auction.SnipeExtension != nil {
		if time.Until(*auction.Event.EndTime) < *auction.SnipeRange {
			*auction.Event.EndTime = auction.Event.EndTime.Add(*auction.SnipeExtension)
			u.SuccessMessage(s, i.ChannelID, &discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{{
					Title:       "**Anti-Snipe Activated!**",
					Description: fmt.Sprintf("New End Time: <t:%d>", auction.Event.EndTime.Unix()),
				}},
			})
		}
	}

	//TODO ADD USE CURRENCY

	//Checking if the auction has ended.
	if auction.Event.EndTime.Before(time.Now()) {
		return fmt.Errorf("cannot Bid, Auction has ended")
	}

	// When there is a winning bid
	if bid < auction.Bid {
		return fmt.Errorf("You must bid higher than: " + u.PriceFormat(auction.Bid, auction.Event.GuildID, auction.Currency))
	}

	// When there is no winning bid
	if bid == auction.Bid && auction.WinnerID != nil {
		return fmt.Errorf("You must bid higher than: " + u.PriceFormat(auction.Bid, auction.Event.GuildID, auction.Currency))
	}

	//TODO Refactor some of this Potentially add Buyout Button instead.
	if auction.Buyout == nil || bid < *auction.Buyout {

		//Checking if the auction is capped and the current winner is bidding.
		if (auction.WinnerID != nil) && (auction.IncrementMax != nil) && (i.Member.User.ID == *auction.WinnerID) {
			return fmt.Errorf("cannot out bid yourself on a capped bid auction")
		}

		//Checking if integer only bidding is enabled.
		if auction.IntegerOnly && bid != math.Floor(bid) {
			return fmt.Errorf("Your bid must be an integer for this auction! For example: " + fmt.Sprint(math.Floor(bid)) + " instead of " + u.PriceFormat(bid, auction.Event.GuildID, auction.Currency))
		}

		//Checking if bid is higher than minimum increment.
		if auction.IncrementMin != nil && bid-auction.Bid < *auction.IncrementMin {
			return fmt.Errorf("Bid must be higher than the previous bid by: %s\n\u200b", u.PriceFormat(*auction.IncrementMin, auction.Event.GuildID, auction.Currency))
		}

		//Checking if bid is lower than maximum increment.
		if auction.IncrementMax != nil && bid-auction.Bid > *auction.IncrementMax {
			return fmt.Errorf("Bid must be no more than %s higher than the previous bid. \n\u200b", u.PriceFormat(*auction.IncrementMax, auction.Event.GuildID, auction.Currency))
		}
	}

	auction.Bid = bid
	newBid := fmt.Sprintf("-> %s: %s", i.Member.User.Username, u.PriceFormat(bid, auction.Event.GuildID, auction.Currency))

	if auction.BidHistory == nil {
		auction.BidHistory = u.Ptr(newBid)
	} else {
		auction.BidHistory = u.Ptr(fmt.Sprintf("%s\n%s", *auction.BidHistory, newBid))
	}

	if auction.Buyout != nil && bid >= *auction.Buyout {
		auction.Event.EndTime = u.Ptr(time.Now())
	}

	err = database.SaveAuction(&auction)
	if err != nil {
		return err
	}

	//TODO Handle setting snipe range and snipe extension mid auction

	m := AuctionMessageFormat(auction)

	_, err = u.SuccessMessageEdit(s, i.ChannelID, *auction.Event.MessageID, &m)
	if err != nil {
		return err
	}

	return u.SuccessResponse(s, i, &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Bid has been successfully placed!",
			},
		},
	})
}

func AuctionEnd(s *discordgo.Session, channelID, guildID string) error {

	auction, err := database.GetAuctionData(channelID)
	if err != nil {
		return err
	}

	auctionSettings, err := database.GetAuctionSettings(guildID)
	if err != nil {
		return err
	}

	//TODO Add in a check for if the auction has ended already

	//TODO Handle the case when an auction message is not fetchable.
	/* message, err := s.ChannelMessage(channelID, *auction.Event.MessageID)
	if err != nil {
		return err
	} */

	if auctionSettings.LogChannel == nil {
		fmt.Println("Log channel has not been set for guild: " + guildID)
		// TODO Instead of sending an error message, update the auction with the error
		_, err := u.ErrorMessage(s, channelID, fmt.Errorf("Auction cannot end because log channel has not been set. Please setup an auction log using `/settings auction`. You might need to end the auction manually after setting the channel."))
		return err
	}





	//TODO Undisable claiming lmao
	err = ClaimOutput(s, auctionMap, "Auction")
	if err != nil {
		return fmt.Errorf("Claim Output Error: " + err.Error())
	} 

	message := AuctionMessageFormat(auction)

	//TODO Do something silly and wacky when the auction ends again
	/* 	if message != nil {
		message.Embeds = append(messageEmbeds.Embeds, &discordgo.MessageEmbed{
			Title:       "Auction has ended!",
			Description: "Thank you for participating!",
			Color:       0x32CD32,
			Image: &discordgo.MessageEmbedImage{
				URL: "https://c.tenor.com/MvFFZxXwdpwAAAAC/sold-ray.gif",
			},
		})
		message.Components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Support Server",
						Style: discordgo.LinkButton,
						Emoji: discordgo.ComponentEmoji{
							Name:     "logo",
							ID:       "889025400120950804",
							Animated: false,
						},
						URL: "https://discord.gg/RxP2z5NGtj",
					},
				},
			},
		}
		_, err = s.ChannelMessageEditComplex(message)
		if err != nil {
			fmt.Println(err)
		}
	}*/

	//TODO Add functionality for channel lock again
	/* 	if auction != true {
		time.Sleep(30 * time.Second)
		_, err = s.ChannelDelete(channelID)
		if err != nil {
			fmt.Println(err)
		}
	} */

	return nil
}

func AuctionEndButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	auction, err := database.GetAuctionData(i.ChannelID)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label: "Delete Auction Channel",
								Style: 4,
								Emoji: discordgo.ComponentEmoji{
									Name: "🛑",
								},
								CustomID: "delete_auction_channel",
							},
							discordgo.Button{
								Label: "Support Server",
								Style: discordgo.LinkButton,
								Emoji: discordgo.ComponentEmoji{
									Name:     "logo",
									ID:       "889025400120950804",
									Animated: false,
								},
								URL: "https://discord.gg/RxP2z5NGtj",
							},
						},
					},
				},
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Error",
						Description: "There was an error. Please contact support if you need help.\n**If you would like to close this auction channel, press the button below.**",
						Color:       0xff0000,
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "**Error Message:**",
								Value: err.Error(),
							},
						},
					},
				},
				Flags: 64,
			},
		})
	}

	//TODO Configure this to allow setting different perms or do a poll about manage server.
	if (i.Member.Permissions&discordgo.PermissionAdministrator == 0) && (i.Member.User.ID != auction.Event.Host) {
		return fmt.Errorf("User must have administrator permissions or be host to run this command")
	}

	err = u.SuccessResponse(s, i, &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{{
			Title:       "End Auction",
			Description: "Please wait while the auction is ending...",
		}},
	})
	if err != nil {
		return err
	}

	auction.Event.EndTime = u.Ptr(time.Now())

	err = database.SaveAuction(&auction)
	if err != nil {
		return err
	}

	err = AuctionEnd(s, i.ChannelID, i.GuildID)
	if err != nil {
		return err
	}
	return nil
}

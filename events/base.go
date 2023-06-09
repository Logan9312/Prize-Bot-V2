package events

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	. "gitlab.com/logan9312/discord-auction-bot/helpers"
)

type Event struct {
	ID          uint `gorm:"primaryKey"`
	BotID       string
	GuildID     string
	Host        string
	Item        string
	ChannelID   *string
	MessageID   *string
	StartTime   *time.Time
	EndTime     *time.Time
	ImageURL    *string
	Description *string
	Note        *string
}

func (event Event) StartTimers() {

}

func NewEvent(botID string, i *discordgo.InteractionCreate, options map[string]any) (Event, error) {
	eventData := Event{
		BotID:   botID,
		GuildID: i.GuildID,
		Host:    i.Member.User.ID,
		Item:    options["item"].(string),
	}

	if options["duration"] != nil {
		duration, err := ParseTime(options["duration"].(string))
		if err != nil {
			return Event{}, fmt.Errorf("Error parsing time input: %w", err)
		}
		eventData.EndTime = Ptr(time.Now().Add(duration))
	}

	if options["image"] != nil {
		eventData.ImageURL = ImageToURL(i, options["image"].(string))
	}

	duration, err := ParseTime(options["duration"].(string))
	if err != nil {
		return eventData, fmt.Errorf("Error parsing time input: %w", err)
	}

	if options["schedule"] != nil {
		startTimeDuration, err := ParseTime(options["schedule"].(string))
		if err != nil {
			return eventData, fmt.Errorf("Error parsing time input: %w", err)
		}
		eventData.StartTime = Ptr(time.Now().Add(startTimeDuration))
		eventData.EndTime = Ptr(eventData.StartTime.Add(duration))
	} else {
		eventData.EndTime = Ptr(time.Now().Add(duration))
	}

	if options["image"] != nil {
		eventData.ImageURL = ImageToURL(i, options["image"].(string))
	}

	if options["description"] != nil {
		eventData.Description = Ptr(options["description"].(string))
	}

	if options["note"] != nil {
		eventData.Note = Ptr(options["note"].(string))
	}

}

// TODO Rework this file
func MessageFormat(data Event) discordgo.MessageSend {

	message := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "",
				Description: fmt.Sprintf("**Host:** <@%s>.\n", data.Host),
				Timestamp:   "",
				Footer:      &discordgo.MessageEmbedFooter{},
				Image:       &discordgo.MessageEmbedImage{},
				Thumbnail:   &discordgo.MessageEmbedThumbnail{},
				Video:       &discordgo.MessageEmbedVideo{},
				Provider:    &discordgo.MessageEmbedProvider{},
				Author:      &discordgo.MessageEmbedAuthor{},
				Fields:      []*discordgo.MessageEmbedField{},
			},
		},
	}

	if data.ImageURL != nil {
		message.Embeds[0].Image = &discordgo.MessageEmbedImage{
			URL: *data.ImageURL,
		}
	}

	description := fmt.Sprintf("**Host:** <@%s>.\n", data.Host)

	if data.Description != nil {
		description += fmt.Sprintf("**Description:** %s\n", data.Description)
	}

	auctionfields := []*discordgo.MessageEmbedField{
		{
			Name:  "__**Details:**__",
			Value: description,
		},
	}

	if data.EndTime != nil {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**End Time**__",
			Value:  fmt.Sprintf("<t:%d:R>", data.EndTime.Unix()),
			Inline: true,
		})
	}

	if data["winner"] != nil {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Current Winner**__",
			Value:  fmt.Sprintf("<@%s>", data["winner"]),
			Inline: true,
		})
	}

	if eventType == EventTypeAuction {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:  "__**How to Bid**__",
			Value: "Use the /bid command or type `/bid {value}` in chat\nEx: `/bid 550`\n",
		})
	}

	if eventType == EventTypeGiveaway {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**How to Enter**__",
			Value:  "To enter, select the 🎁 reaction below! Removing your reaction will remove your entry.",
			Inline: false,
		})
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		fmt.Println("Error fetching guild: ", err)
		return &discordgo.InteractionResponseData{}, err
	}

	if data["alert_role"] != nil {
		content = fmt.Sprintf("<@&%s>", strings.Trim(data["alert_role"].(string), " "))
	}
	if data["alert_role"] == guildID {
		content = "@everyone"
	}

	components := []discordgo.MessageComponent{}

	if eventType == EventTypeAuction {
		components = []discordgo.MessageComponent{
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
	}

	if eventType == EventTypeShop {
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Purchase",
						Style:    discordgo.SuccessButton,
						CustomID: "additem",
					},
				},
			},
		}
	}

	if data["bid_history"] != nil {
		if len(data["bid_history"].(string)) > 4095 {
			data["bid_history"] = data["bid_history"].(string)[len(data["bid_history"].(string))-4095:]
		}
		embeds = []*discordgo.MessageEmbed{{
			Title:       "**Bid History**",
			Description: data["bid_history"].(string),
			Color:       0x8073ff,
			Image: &discordgo.MessageEmbedImage{
				URL: "https://i.imgur.com/9wo7diC.png",
			},
		}}
	}

	return &discordgo.InteractionResponseData{
		Content:    content,
		Title:      fmt.Sprintf("%s Item: __**%s**__", eventType, data["item"]),
		Fields:     auctionfields,
		Thumbnail:  &discordgo.MessageEmbedThumbnail{URL: guild.IconURL()},
		Image:      &discordgo.MessageEmbedImage{URL: imageURL},
		Components: components,
		Embeds:     embeds,
		Files:      []*discordgo.File{},
	}, nil
}

func HasRole(m *discordgo.Member, roleID string) bool {
	for _, role := range m.Roles {
		if role == roleID {
			return true
		}
	}
	return false
}

package events

import (
	"fmt"
	"time"

	"github.com/Logan9312/Prize-Bot-V2/database"
	u "github.com/Logan9312/Prize-Bot-V2/utils"
	"github.com/bwmarrin/discordgo"
)

func NewEvent(botID string, i *discordgo.InteractionCreate, options map[string]any) (database.Event, error) {
	eventData := database.Event{
		BotID:   botID,
		GuildID: i.GuildID,
		Host:    i.Member.User.ID,
		Item:    options["item"].(string),
	}

	if options["duration"] != nil {
		duration, err := u.ParseTime(options["duration"].(string))
		if err != nil {
			return database.Event{}, fmt.Errorf("error parsing time input: %w", err)
		}
		eventData.EndTime = u.Ptr(time.Now().Add(duration))
	}

	if options["image"] != nil {
		eventData.ImageURL = u.ImageToURL(i, options["image"].(string))
	}

	duration, err := u.ParseTime(options["duration"].(string))
	if err != nil {
		return eventData, fmt.Errorf("error parsing time input: %w", err)
	}

	if options["schedule"] != nil {
		startTimeDuration, err := u.ParseTime(options["schedule"].(string))
		if err != nil {
			return eventData, fmt.Errorf("error parsing time input: %w", err)
		}
		eventData.StartTime = u.Ptr(time.Now().Add(startTimeDuration))
		eventData.EndTime = u.Ptr(eventData.StartTime.Add(duration))
	} else {
		eventData.EndTime = u.Ptr(time.Now().Add(duration))
	}

	if options["image"] != nil {
		eventData.ImageURL = u.ImageToURL(i, options["image"].(string))
	}

	if options["description"] != nil {
		eventData.Description = u.Ptr(options["description"].(string))
	}

	if options["note"] != nil {
		eventData.Note = u.Ptr(options["note"].(string))
	}

	return eventData, nil
}

// TODO Rework this file

func MessageFormat(data database.Event) discordgo.MessageSend {

	message := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: fmt.Sprintf("Item: __**%s**__", data.Item),
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "__**Details:**__",
						Value: fmt.Sprintf("**Host:** <@%s>.\n", data.Host),
					},
				},
			},
		},
	}

	if data.ImageURL != nil {
		message.Embeds[0].Image = &discordgo.MessageEmbedImage{
			URL: *data.ImageURL,
		}
	}

	if data.Description != nil {
		message.Embeds[0].Fields[0].Value += fmt.Sprintf("**Description:** %s\n", *data.Description)
	}

	if data.EndTime != nil {
		message.Embeds[0].Fields = append(message.Embeds[0].Fields, &discordgo.MessageEmbedField{
			Name:   "__**End Time**__",
			Value:  fmt.Sprintf("<t:%d:R>", data.EndTime.Unix()),
			Inline: true,
		})
	}

	//TODO Copy this over to giveaway format
	/*if eventType == EventTypeGiveaway {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**How to Enter**__",
			Value:  "To enter, select the 🎁 reaction below! Removing your reaction will remove your entry.",
			Inline: false,
		})
	}*/

	if data.AlertRole == &data.GuildID {
		message.Content = "@everyone"
	} else if data.AlertRole != nil {
		message.Content = fmt.Sprintf("<@&%s>", *data.AlertRole)
	}

	//TODO Add to shop format
	/*if eventType == EventTypeShop {
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
	}*/

	return message
}

func HasRole(m *discordgo.Member, roleID string) bool {
	for _, role := range m.Roles {
		if role == roleID {
			return true
		}
	}
	return false
}

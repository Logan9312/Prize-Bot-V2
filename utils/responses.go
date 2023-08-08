package utils

import (
	"fmt"
	"runtime"

	"github.com/bwmarrin/discordgo"
)

const (
	ResponseColorRed    = 0xff0000
	ResponseColorPurple = 0x8073ff
	ResponseColorOrange = 0xff6700
	ResponseColorYellow = 0xffd700
)

var SupportButton = discordgo.Button{
	Label: "Support Server",
	Style: discordgo.LinkButton,
	Emoji: discordgo.ComponentEmoji{
		Name:     "logo",
		ID:       "889025400120950804",
		Animated: false,
	},
	URL: "https://discord.gg/RxP2z5NGtj",
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, resp *discordgo.InteractionResponseData) error {
	resp.Flags = 64
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}

func followupRespond(s *discordgo.Session, i *discordgo.InteractionCreate, resp *discordgo.WebhookParams) (*discordgo.Message, error) {
	resp.Flags = 64
	return s.FollowupMessageCreate(i.Interaction, true, resp)
}


func SuccessResponse(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.InteractionResponseData) error {
	for _, embed := range data.Embeds {
		if embed.Color == 0 {
			embed.Color = ResponseColorPurple
		}
	}
	return respond(s, i, data)
}

func FollowUpSuccessResponse(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.WebhookParams) (*discordgo.Message, error) {
	for _, embed := range data.Embeds {
		if embed.Color == 0 {
			embed.Color = ResponseColorPurple
		}
	}
	return followupRespond(s, i, data)
}

func SuccessMessage(s *discordgo.Session, channelID string, data *discordgo.MessageSend) (*discordgo.Message, error) {

	for _, embed := range data.Embeds {
		if embed.Color == 0 {
			embed.Color = ResponseColorPurple
		}
	}
	return s.ChannelMessageSendComplex(channelID, data)
}

func SuccessMessageEdit(s *discordgo.Session, channelID, messageID string, data *discordgo.MessageSend) (*discordgo.Message, error) {
	for _, embed := range data.Embeds {
		if embed.Color == 0 {
			embed.Color = ResponseColorPurple
		}
	}
	return s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content:         &data.Content,
		Components:      data.Components,
		AllowedMentions: data.AllowedMentions,
		Embeds:          data.Embeds,
		Files:           data.Files,
		ID:              messageID,
		Channel:         channelID,
		Embed:           data.Embed,
	})
}

func ExperimentalResponse(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.InteractionResponseData) error {
	if len(data.Embeds) == 0 {
		return Errorfmt(fmt.Errorf("no embeds provided"))
	}
	data.Embeds = append(data.Embeds, &discordgo.MessageEmbed{
		Title:       "**Experimental Feature!**",
		Description: "This feature is highly experimental, and for that reason may not function perfectly as intended. Often problems with this command are more frequent in large servers.",
	})

	for _, embed := range data.Embeds {
		if embed.Color == 0 {
			embed.Color = ResponseColorOrange
		}
	}
	if data.Components == nil {
		data.Components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					SupportButton,
				},
			},
		}
	} else {
		data.Components = append(data.Components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				SupportButton,
			},
		})
	}

	return respond(s, i, data)
}

func Errorfmt(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("%s:%d %w", file, line, err)
}

func ErrorComponentfmt(err error) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				SupportButton,
			},
		},
	}
}

func ErrorEmbedfmt(err error) []*discordgo.MessageEmbed {
	return []*discordgo.MessageEmbed{
		{
			Title:       "Error",
			Description: "There was an error. Please contact support if you need help.",
			Color:       ResponseColorRed,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "**Error Message:**",
					Value: err.Error(),
				},
				{
					Name:  "**Stack Trace:**",
				},
			},
		},
	}
}

func ErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) error {

	fmt.Println(Errorfmt(err).Error())

	return respond(s, i, &discordgo.InteractionResponseData{
		Components: ErrorComponentfmt(err),
		Embeds:     ErrorEmbedfmt(err),
		Flags:      64,
	})
}

func FollowUpErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) (*discordgo.Message, error) {

	fmt.Println(Errorfmt(err).Error())

	return followupRespond(s, i, &discordgo.WebhookParams{
		Components: ErrorComponentfmt(err),
		Embeds:     ErrorEmbedfmt(err),
	})
}

func ErrorMessage(s *discordgo.Session, channelID string, err error) (*discordgo.Message, error) {

	fmt.Println(Errorfmt(err).Error())

	return s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     ErrorEmbedfmt(err),
		Components: ErrorComponentfmt(err),
	})
}

func PremiumError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {

	return respond(s, i, &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "__**Premium Users Only**__",
				Description: "This function is restricted to premium users. Please check out `/premium info` for more information.",
				Color:       ResponseColorYellow,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "**Message**",
						Value:  message,
						Inline: false,
					},
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://i.imgur.com/SjJT2xp.png",
				},
			},
		},
	})
}

package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Logan9312/Prize-Bot-V2/database"
	u "github.com/Logan9312/Prize-Bot-V2/utils"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//TODO Rework this whole file

var (
	ClaimingCommandOption = &discordgo.ApplicationCommandOption	{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "claiming",
		Description: "Edit your claiming settings",
		Options: []*discordgo.ApplicationCommandOption{
			CommandOptionCategory,
			CommandOptionLogChannel,
			CommandOptionInstructions,
			CommandOptionDisableClaiming,
			//TODO Make channel prefix work for claiming
			//CommandOptionChannelPrefix,
		},
	}

	CommandOptionAuctions = &discordgo.ApplicationCommandOption	{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "auctions",
		Description: "Edit your auction settings",
		Options: []*discordgo.ApplicationCommandOption{
			CommandOptionLogChannel,
			CommandOptionCategory,
			CommandOptionAlertRole,
			CommandOptionIntegerOnly,
			CommandOptionUseCurrency,
			CommandOptionSnipeExtension,
			CommandOptionSnipeRange,
			CommandOptionChannelLock,
			CommandOptionChannelPrefix,
		},
	}
)

var SettingsCommand = discordgo.ApplicationCommand{
	Name:                     "settings",
	Description:              "Change any settings for your bot",
	DefaultMemberPermissions: u.Ptr(int64(discordgo.PermissionManageServer)),
	DMPermission:             new(bool),
	Options: []*discordgo.ApplicationCommandOption{
		ClaimingCommandOption,

		{
			Type:         discordgo.ApplicationCommandOptionSubCommand,
			Name:         "giveaways",
			Description:  "Edit your giveaway settings",
			ChannelTypes: []discordgo.ChannelType{},
			Required:     false,
			Options: []*discordgo.ApplicationCommandOption{
				CommandOptionAlertRole,
				CommandOptionLogChannel,
				CommandOptionAutoEnter,
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "currency",
			Description: "Edit your currency settings",
			Options: []*discordgo.ApplicationCommandOption{
				//TODO Rework for multiple currencies
				CommandOptionCurrency,
				CommandOptionCurrencySide,
			},
		},
		/*{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "shop",
			Description: "Edit your shop settings",
			Options: []*discordgo.ApplicationCommandOption{},
		},*/
	},
}

func Settings(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := u.ParseSubCommand(i)
	responseFields := []*discordgo.MessageEmbedField{}
	menuOptions := []discordgo.SelectMenuOption{}
	var commandData discordgo.ApplicationCommandOption
	var antiSnipeDescription string
	var model interface{}
	var customID string
	var content string
	var err error

	switch i.ApplicationCommandData().Options[0].Name {
	case "claiming":
		model = &database.ClaimSetup{
			GuildID: i.GuildID,
		}
		commandData = *SettingsCommand.Options[0]
		customID = "clear_claim_setup"
	case "auctions":
		model = &database.AuctionSetup{
			GuildID: i.GuildID,
		}
		commandData = *SettingsCommand.Options[1]
		customID = "clear_auction_setup"
	case "giveaways":
		model = &database.GiveawaySetup{
			GuildID: i.GuildID,
		}
		commandData = *SettingsCommand.Options[2]
		customID = "clear_giveaway_setup"
	case "currency":
		model = &database.CurrencySetup{
			GuildID: i.GuildID,
		}
		commandData = *SettingsCommand.Options[3]
		customID = "clear_currency_setup"
	case "shop":
		model = &database.ShopSetup{
			GuildID: i.GuildID,
		}
		commandData = *SettingsCommand.Options[4]
		customID = "clear_shop_setup"
	default:
		return fmt.Errorf("unknown settings type: %s", i.ApplicationCommandData().Options[0].Name)
	}

	options["guild_id"] = i.GuildID

	for key := range options {
		content += fmt.Sprintf("• %s has been successfully set.\n", cases.Title(language.English).String(strings.ReplaceAll(key, "_", " ")))

		switch key {
		case "snipe_extension", "snipe_range":
			options[key], err = u.ParseTime(options[key].(string))
		}
		if err != nil {
			return err
		}
	}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(model).Create(options)
	if result.Error != nil {
		return result.Error
	}

	result = database.DB.Model(model).Updates(options)
	if result.Error != nil {
		return result.Error
	}

	result = database.DB.Model(model).First(options, i.GuildID)
	if result.Error != nil {
		return result.Error
	}

	for _, v := range commandData.Options {
		if !strings.Contains(v.Name, "snipe") {
			switch {
			case options[v.Name] == nil || options[v.Name] == "":
				options[v.Name] = "Not Set"
			case fmt.Sprintf("%T", options[v.Name]) == "bool":
				options[v.Name] = strings.Title(fmt.Sprint(options[v.Name]))
			case strings.Contains(v.Name, "role"):
				options[v.Name] = fmt.Sprintf("<@&%s>", options[v.Name])
			case strings.Contains(v.Name, "channel") && v.Name != "channel_prefix":
				options[v.Name] = fmt.Sprintf("<#%s>", options[v.Name])
			case strings.Contains(v.Name, "category"):
				category, err := s.Channel(options[v.Name].(string))
				if err != nil {
					fmt.Println("Category Error:", err)
					options[v.Name] = "Error Displaying Category: " + err.Error()
				} else {
					options[v.Name] = category.Name
				}
			}
			responseFields = append(responseFields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("**%s**", strings.Title(strings.ReplaceAll(v.Name, "_", " "))),
				Value: fmt.Sprint(options[v.Name]),
			})
		}
		menuOptions = append(menuOptions, discordgo.SelectMenuOption{
			Label:       strings.Title(strings.ReplaceAll(v.Name, "_", " ")),
			Value:       v.Name,
			Description: v.Description,
		})
	}

	if i.ApplicationCommandData().Options[0].Name == "auctions" {
		if options["snipe_range"] != nil && options["snipe_extension"] != nil {
			antiSnipeDescription = fmt.Sprintf("If a bid is placed within %s of the auction ending, it will be extended by %s.", options["snipe_range"].(time.Duration).String(), options["snipe_extension"].(time.Duration).String())
		} else {
			antiSnipeDescription = "Anti Snipe Disabled. To enable, set both snipe_extension and snipe_range"
		}

		responseFields = append(responseFields, &discordgo.MessageEmbedField{
			Name:  "**Anti Snipe**",
			Value: antiSnipeDescription,
		})
	}
	return u.SuccessResponse(s, i, &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       fmt.Sprintf("%s Setup", strings.Title(i.ApplicationCommandData().Name)),
				Description: content,
				Fields:      responseFields,
			},
		},

		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    customID,
						Placeholder: "Clear Setup Options",
						MinValues:   u.Ptr(1),
						MaxValues:   len(commandData.Options),
						Options:     menuOptions,
					},
				},
			},
		},
	})
}

func CurrencySetupClearButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := i.MessageComponentData().Values
	clearedMap := map[string]interface{}{}

	info := database.CurrencySetup{
		GuildID: i.GuildID,
	}

	clearedSettings := "No Settings Cleared!"
	if len(options) > 0 {
		clearedSettings = ""
	}

	for _, v := range options {
		clearedSettings += fmt.Sprintf("• %s\n", strings.Title(strings.ReplaceAll(v, "_", " ")))
		clearedMap[v] = gorm.Expr("NULL")
	}

	database.DB.Model(&info).Updates(clearedMap)

	return u.SuccessResponse(s, i, &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "**Cleared Auction Settings**",
				Description: "You have successfully cleared the following settings. Run `/settings auction` to see your changes.",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "**Cleared Settings**",
						Value: clearedSettings,
					},
				},
			},
		},
	})
}

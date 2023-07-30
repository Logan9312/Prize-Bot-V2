package connect

import (
	"fmt"

	c "github.com/Logan9312/Prize-Bot-V2/commands"
	"github.com/bwmarrin/discordgo"
)

type slashCommands struct {
	Local, Prod, Dev []*discordgo.ApplicationCommand
}

var BotCommands = slashCommands{
	Local: []*discordgo.ApplicationCommand{
		//&c.ShopCommand,
		//&c.QuestCommand,
		//&QuestContextMenu,
	},
	Prod: []*discordgo.ApplicationCommand{
		//&c.ProfileCommand,
		&c.AuctionCommand,
		&c.BidCommand,
		&c.GiveawayCommand,
		//&c.PrivacyCommand,
		&c.ClaimCommand,
		//&c.PremiumCommand,
		&c.SettingsCommand,
		&c.CurrencyCommand,
		&WhitelabelCommand,
	},
	Dev: []*discordgo.ApplicationCommand{
		//&c.DevCommand,
		//&c.StatsCommand,
	},
}

func BotConnect(token, environment string) (*discordgo.Session, error) {

	BotCommands.Local = append(BotCommands.Local, BotCommands.Prod...)
	BotCommands.Local = append(BotCommands.Local, BotCommands.Dev...)

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return s, fmt.Errorf("Discordgo.New Error: %w", err)
	}

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages

	RegisterHandlers(s)

	err = s.Open()
	if err != nil {
		return s, fmt.Errorf("Failed to open a websocket connection with discord. Likely due to an invalid token. %w", err)
	}

	fmt.Println(s.State.User.Username, " Starting Up...")

	//Clears all guild commands
	//s.LogLevel = discordgo.LogInformational
	for _, v := range s.State.Guilds {
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, []*discordgo.ApplicationCommand{})
		fmt.Println("Commands added to guild: " + v.Name)
		if err != nil {
			return s, fmt.Errorf("bulk Overwrite Local Command Error: %w", err)
		}
	}

	//Builds prod commands
	commands := BotCommands.Prod
	if environment == "local" {
		commands = BotCommands.Local
	}

	//Builds global commands
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", commands)
	if err != nil {
		return s, fmt.Errorf("bulk Overwrite Prod Command Error: %w", err)
	}

	//Builds dev commands
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "915767892467920967", BotCommands.Dev)
	if err != nil {
		return s, fmt.Errorf("Bulk Overwrite Dev Command Error: %w", err)
	}

	fmt.Println(s.State.User.Username + " bot startup complete!")

	return s, nil
}

//TODO Startup timers

package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ParseCommand(i *discordgo.InteractionCreate) map[string]interface{} {
	data := i.ApplicationCommandData().Options
	options := make(map[string]interface{})
	parseOptions(data, options)
	if options["image"] != nil {
		options["image"] = ImageToURL(i, options["image"].(string))
	}
	return options
}

func parseOptions(data []*discordgo.ApplicationCommandInteractionDataOption, options map[string]interface{}) {
	for _, option := range data {
		if option.Options == nil {
			options[option.Name] = option.Value
		} else {
			subOptions := make(map[string]interface{})
			parseOptions(option.Options, subOptions)
			options[option.Name] = subOptions
		}
	}
}

func ParseTime(inputDuration string) (time.Duration, error) {
	if strings.HasSuffix(strings.ToLower(inputDuration), "d") {
		inputDuration = strings.TrimSuffix(strings.ToLower(inputDuration), "d")
		float, err := strconv.ParseFloat(inputDuration, 64)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		inputDuration = fmt.Sprint(float*24) + "h"
	}

	return time.ParseDuration(strings.ToLower(inputDuration))
}

func Ptr[T any](v T) *T {
	return &v
}

func ImageToURL(i *discordgo.InteractionCreate, image string) *string {
	return &i.ApplicationCommandData().Resolved.Attachments[image].URL
}

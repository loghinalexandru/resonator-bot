package command

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/rest"
)

type swearData struct {
	Swear string `json:"swear"`
	Lang  string `json:"lang"`
}

func NewSwear(ctx BotContext) *rest.REST[swearData] {
	url := ctx.SwearsApiURL.String() + "/api/random?lang=%v"
	def := &discordgo.ApplicationCommand{
		Name:        "swear",
		Description: "This command is used to play a TTS message of a random swear!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "language",
				Description: "Specify in which language your swear will be!",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Romanian",
						Value: "ro",
					},
					{
						Name:  "English",
						Value: "en",
					},
					{
						Name:  "French",
						Value: "fr",
					},
				},
			},
		},
	}

	result, err := rest.New(def, url, rest.WithFormatter[swearData](swearFormatter))

	if err != nil {
		panic(err)
	}

	return result
}

func swearFormatter(resp swearData) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("> \"**%s**\"", resp.Swear))
	return sb.String()
}

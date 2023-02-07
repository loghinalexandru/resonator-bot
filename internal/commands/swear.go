package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/rest"
)

type SwearData struct {
	Swear string `json:"swear"`
	Lang  string `json:"lang"`
}

func NewSwear(swearsURL string) *rest.REST {
	out := rest.New(&discordgo.ApplicationCommand{
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
	})

	out.URL = swearsURL + "/api/random?lang=%v"
	out.Type = &SwearData{}
	out.Formatter = swearFormatter

	return &out
}

func swearFormatter(content any) string {
	var sb strings.Builder
	resp, ok := content.(*SwearData)

	if !ok {
		return "Something went wrong!"
	}

	sb.WriteString(fmt.Sprintf("> \"**%s**\"", resp.Swear))
	return sb.String()
}

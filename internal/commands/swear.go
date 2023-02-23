package commands

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/rest"
)

type SwearData struct {
	Swear string `json:"swear"`
	Lang  string `json:"lang"`
}

func NewSwear(swearsURL string, client *http.Client) *rest.REST[SwearData] {
	url := swearsURL + "/api/random?lang=%v"
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

	return rest.New(def, url, client, swearFormatter)
}

func swearFormatter(resp SwearData) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("> \"**%s**\"", resp.Swear))
	return sb.String()
}

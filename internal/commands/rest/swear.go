package rest

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type SwearData struct {
	Swear string `json:"swear"`
	Lang  string `json:"lang"`
}

func NewSwear() *REST {
	out := REST{
		URL:       "http://localhost:3000/api/random?lang=%v",
		TTS:       true,
		Type:      &SwearData{},
		Formatter: swearFormatter,
		Def: &discordgo.ApplicationCommand{
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
		},
	}
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

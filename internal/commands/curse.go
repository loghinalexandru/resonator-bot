package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/playback"
)

func NewCurse(sync *sync.Map, baseURL string) *playback.Playback {
	return playback.New(sync, &discordgo.ApplicationCommand{
		Name:        "curse",
		Description: "This command is used to play a friendly encouragement!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "language",
				Description: "Specify in which language your encouragement will be!",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Romanian",
						Value: "/api/random/file?lang=ro&opus=true",
					},
					{
						Name:  "French",
						Value: "/api/random/file?lang=fr&opus=true",
					},
					{
						Name:  "English",
						Value: "/api/random/file?lang=en&opus=true",
					},
				},
				Required: true,
			},
		},
	},
		playback.WithURL(baseURL),
	)
}

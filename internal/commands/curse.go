package commands

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/playback"
)

func NewCurse(sync *sync.Map, swearsURL string) *playback.Playback {
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
						Value: fmt.Sprintf("%v/api/random/file?lang=ro&opus=true", swearsURL),
					},
					{
						Name:  "French",
						Value: fmt.Sprintf("%v/api/random/file?lang=fr&opus=true", swearsURL),
					},
					{
						Name:  "English",
						Value: fmt.Sprintf("%v/api/random/file?lang=en&opus=true", swearsURL),
					},
				},
				Required: true,
			},
		},
	})
}

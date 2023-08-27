package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/audio"
	"github.com/loghinalexandru/resonator/pkg/playback"
)

const (
	randomPath = "/api/random/file?codec=opus&lang=%v"
)

func NewCurse(ctx BotContext) *playback.Playback {
	result, err := playback.New(ctx.Sync, &discordgo.ApplicationCommand{
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
						Value: "ro",
					},
					{
						Name:  "French",
						Value: "fr",
					},
					{
						Name:  "English",
						Value: "en",
					},
				},
				Required: true,
			},
		},
	},
		playback.WithSource(audio.NewHTTP(ctx.SwearsApiURL.String()+randomPath)),
	)

	if err != nil {
		panic(err)
	}

	return result
}

package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/audio"
	"github.com/loghinalexandru/resonator/pkg/playback"
)

const (
	remotePath = "/api/remote?codec=opus&id=%v"
)

func NewFeed(ctx BotContext) *playback.Playback {
	result, err := playback.New(ctx.Sync, &discordgo.ApplicationCommand{
		Name:        "feed",
		Description: "This command is used to play remote youtube sound!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "Specify youtube id to feed into the bot",
				Required:    true,
			},
		},
	},
		playback.WithSource(audio.NewHTTP(ctx.SwearsApiURL.String()+remotePath)),
	)

	if err != nil {
		panic(err)
	}

	return result
}

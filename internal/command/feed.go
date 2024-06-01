package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/bot"
	"github.com/loghinalexandru/resonator/pkg/playback"
	"github.com/loghinalexandru/resonator/pkg/provider"
)

const (
	remotePath = "/api/remote?codec=opus&id=%v"
)

func newFeed(ctx *bot.Context) *playback.Playback {
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
		playback.WithSource(provider.NewHTTP(ctx.SwearsAPI.String()+remotePath)),
	)

	if err != nil {
		ctx.Logger.Error("Error creating /feed command", "err", err)
	}

	return result
}

package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/playback"
)

const (
	remotePath = "/api/remote?encoding=opus&id=%v"
)

func NewFeed(sync *sync.Map, baseURL string) *playback.Playback {
	return playback.New(sync, &discordgo.ApplicationCommand{
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
		//Fix this to be more easy to extend and not concat strings
		playback.WithURL(baseURL+remotePath),
	)
}

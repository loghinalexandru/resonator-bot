package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/playback"
)

func NewPlay(sync *sync.Map) *playback.Playback {
	return playback.New(sync, &discordgo.ApplicationCommand{
		Name:        "play",
		Description: "This command is used to play a sound in the chat!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "type",
				Description: "Sound type to be played!",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Ara-Ara",
						Value: "misc/ara.dca",
					},
					{
						Name:  "Yoooooooouuu",
						Value: "misc/yoo.dca",
					},
					{
						Name:  "FBI Open Up",
						Value: "misc/fbi.dca",
					},
					{
						Name:  "Hehe Boy!",
						Value: "misc/heheboy.dca",
					},
					{
						Name:  "Yamete Kudasai!",
						Value: "misc/yamete.dca",
					},
					{
						Name:  "No God Please No!",
						Value: "misc/nogod.dca",
					},
					{
						Name:  "Mission failed.",
						Value: "misc/failed.dca",
					},
					{
						Name:  "Death",
						Value: "misc/death.dca",
					},
					{
						Name:  "UwU",
						Value: "misc/uwu.dca",
					},
					{
						Name:  "Fast AF",
						Value: "misc/fastaf.dca",
					},
				},
				Required: true,
			},
		},
	})
}

package playback

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

func NewPlay(sync *sync.Map) *Playback {
	out := Playback{
		Storage: sync,
		Def: &discordgo.ApplicationCommand{
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
					},
					Required: true,
				},
			},
		},
	}
	return &out
}

package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

func playCommand(sync *sync.Map) *types.Playback {
	out := types.Playback{
		Storage: sync,
	}

	result := out.Definition()
	result.Name = "play"
	result.Description = "This command is used to play a sound in the chat!"
	result.Options = append(result.Options, &discordgo.ApplicationCommandOption{
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
	})

	return &out
}

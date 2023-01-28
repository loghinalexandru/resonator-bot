package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

func playCommand() *types.Playback {
	var out types.Playback

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
				Value: "ara.dca",
			},
			{
				Name:  "Yoooooooouuu",
				Value: "yoo.dca",
			},
			{
				Name:  "FBI Open Up",
				Value: "fbi.dca",
			},
			{
				Name:  "Hehe Boy!",
				Value: "heheboy.dca",
			},
			{
				Name:  "Yamete Kudasai!",
				Value: "yamete.dca",
			},
			{
				Name:  "No God Please No!",
				Value: "nogod.dca",
			},
			{
				Name:  "Mission failed.",
				Value: "failed.dca",
			},
			{
				Name:  "Death",
				Value: "death.dca",
			},
		},
		Required: true,
	})

	return &out
}

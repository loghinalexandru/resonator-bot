package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

func reactCommand() *types.Playback {
	result := discordgo.ApplicationCommand{}
	result.Name = "react"
	result.Description = "This command is used to react with a sound in the chat!"
	result.Options = append(result.Options, &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "reaction",
		Description: "Reaction to be played!",
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Eh?",
				Value: "eh.dca",
			},
			{
				Name:  "Alta intrebare.",
				Value: "intrebare.dca",
			},
			{
				Name:  "Yass",
				Value: "yass.dca",
			},
		},
		Required: true,
	})

	return &types.Playback{
		Def: result,
	}
}

package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

func NewReact(sync *sync.Map) *types.Playback {
	out := types.Playback{
		Storage: sync,
		Def: &discordgo.ApplicationCommand{
			Name:        "react",
			Description: "This command is used to react with a sound in the chat!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reaction",
					Description: "Reaction to be played!",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Eh?",
							Value: "misc/eh.dca",
						},
						{
							Name:  "Alta intrebare.",
							Value: "misc/intrebare.dca",
						},
						{
							Name:  "Yass",
							Value: "misc/yass.dca",
						},
						{
							Name:  "Bruuh!",
							Value: "misc/bruh.dca",
						},
						{
							Name:  "Bagmias Pl.",
							Value: "misc/pl.dca",
						},
						{
							Name:  "To be continued...",
							Value: "misc/continued.dca",
						},
						{
							Name:  "Directed By Robert B. Weide",
							Value: "misc/directedby.dca",
						},
					},
					Required: true,
				},
			},
		},
	}

	return &out
}

package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

func NewAnime() *types.Kitsu {
	out := types.Kitsu{
		URL: "https://kitsu.io/api/edge/anime?filter[text]=%v",
		Def: &discordgo.ApplicationCommand{
			Name:        "anime",
			Description: "This command is used find anime via kitsu API!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keyword",
					Description: "Keyword to search for",
					Required:    true,
				},
			},
		},
	}
	return &out
}

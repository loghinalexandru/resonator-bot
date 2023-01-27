package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

func playCommand() *types.Playback {
	result := discordgo.ApplicationCommand{}
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
				Name:  "Hai mai repede!",
				Value: "repede.dca",
			},
			{
				Name:  "Yoooooooouuu",
				Value: "yoo.dca",
			},
			{
				Name:  "Bagmias Pl.",
				Value: "pl.dca",
			},
			{
				Name:  "FBI Open Up",
				Value: "fbi.dca",
			},
			{
				Name:  "Fa nu mai vomita!",
				Value: "vomita.dca",
			},
			{
				Name:  "To be continued...",
				Value: "continued.dca",
			},
			{
				Name:  "Il bag in ma-sa!",
				Value: "baginmasa.dca",
			},
			{
				Name:  "Da? Vrei ceas?",
				Value: "muienuvrei.dca",
			},
			{
				Name:  "Dau flash!",
				Value: "flash.dca",
			},
			{
				Name:  "Bruuh!",
				Value: "bruh.dca",
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
				Name:  "Directed By Robert B. Weide",
				Value: "directedby.dca",
			},
			{
				Name:  "Sarut-mana",
				Value: "sarutmana.dca",
			},
			{
				Name:  "No God Please No!",
				Value: "nogod.dca",
			},
			{
				Name:  "La culcare!",
				Value: "laculcare.dca",
			},
			{
				Name:  "Da tu cu stomacul ce ai?",
				Value: "stomacul.dca",
			},
			{
				Name:  "Ma tu carti citesti?",
				Value: "carticitesti.dca",
			},
			{
				Name:  "N-am facut asta niciodata!",
				Value: "narerost.dca",
			},
			{
				Name:  "Paul, vieni qui.",
				Value: "sanfranciscu.dca",
			},
			{
				Name:  "Prin puterea zeilor!",
				Value: "putereazeilor.dca",
			},
			{
				Name:  "Mission failed.",
				Value: "failed.dca",
			},
			{
				Name:  "Death",
				Value: "death.dca",
			},
			{
				Name:  "Eh?",
				Value: "eh.dca",
			},
			// {
			// 	Name:  "Alta intrebare.",
			// 	Value: "intrebare.dca",
			// },
			// {
			// 	Name:  "Yass",
			// 	Value: "yass.dca",
			// },
		},
		Required: true,
	})

	return &types.Playback{
		Def: result,
	}
}

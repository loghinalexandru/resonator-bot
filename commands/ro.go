package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

func roCommand() *types.Playback {
	var out types.Playback

	result := out.Definition()
	result.Name = "ro"
	result.Description = "This command is used to play a romanian sound in the chat!"
	result.Options = append(result.Options, &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "type",
		Description: "Sound type to be played!",
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Hai mai repede!",
				Value: "repede.dca",
			},
			{
				Name:  "Fa nu mai vomita!",
				Value: "vomita.dca",
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
				Name:  "Sarut-mana",
				Value: "sarutmana.dca",
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
		},
		Required: true,
	})

	return &out
}

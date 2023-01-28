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
				Value: "misc/repede.dca",
			},
			{
				Name:  "Fa nu mai vomita!",
				Value: "misc/vomita.dca",
			},
			{
				Name:  "Il bag in ma-sa!",
				Value: "misc/baginmasa.dca",
			},
			{
				Name:  "Da? Vrei ceas?",
				Value: "misc/muienuvrei.dca",
			},
			{
				Name:  "Dau flash!",
				Value: "misc/flash.dca",
			},
			{
				Name:  "Sarut-mana",
				Value: "misc/sarutmana.dca",
			},
			{
				Name:  "La culcare!",
				Value: "misc/laculcare.dca",
			},
			{
				Name:  "Da tu cu stomacul ce ai?",
				Value: "misc/stomacul.dca",
			},
			{
				Name:  "Ma tu carti citesti?",
				Value: "misc/carticitesti.dca",
			},
			{
				Name:  "N-am facut asta niciodata!",
				Value: "misc/narerost.dca",
			},
			{
				Name:  "Paul, vieni qui.",
				Value: "misc/sanfranciscu.dca",
			},
			{
				Name:  "Prin puterea zeilor!",
				Value: "misc/putereazeilor.dca",
			},
		},
		Required: true,
	})

	return &out
}

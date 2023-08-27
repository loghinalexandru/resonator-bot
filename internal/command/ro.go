package command

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/playback"
)

func NewRo(sync *sync.Map) *playback.Playback {
	result, err := playback.New(sync, &discordgo.ApplicationCommand{
		Name:        "ro",
		Description: "This command is used to play a romanian sound in the chat!",
		Options: []*discordgo.ApplicationCommandOption{
			{
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
						Name:  "Japoniaaa!",
						Value: "misc/japonia.dca",
					},
					{
						Name:  "Ma tu carti citesti?",
						Value: "misc/carticitesti.dca",
					},
					{
						Name:  "Pielea pulii...",
						Value: "misc/pieleapulii.dca",
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
					{
						Name:  "Hu sha sh",
						Value: "misc/hushash.dca",
					},
					{
						Name:  "Ceausesculeee!",
						Value: "misc/iesiunpic.dca",
					},
					{
						Name:  "Schema 72",
						Value: "misc/schema.dca",
					},
					{
						Name:  "Taci fa analfabeto!",
						Value: "misc/tacifa.dca",
					},
					{
						Name:  "E de la gaze",
						Value: "misc/delagaze.dca",
					},
					{
						Name:  "Dupa bere",
						Value: "misc/dupabere.dca",
					},
					{
						Name:  "Am baut sveps",
						Value: "misc/sveps.dca",
					},
					{
						Name:  "Vorbiti civilizat",
						Value: "misc/vorbiticivilizat.dca",
					},
					{
						Name:  "O pui pe fata",
						Value: "misc/opuipefata.dca",
					},
				},
				Required: true,
			},
		},
	})

	if err != nil {
		panic(err)
	}

	return result
}

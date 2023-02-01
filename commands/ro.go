package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

func NewRo(sync *sync.Map) *types.Playback {
	out := types.Playback{
		Storage: sync,
		Voice: func(sess *discordgo.Session, guildID, voiceID string, mute, deaf bool) (*discordgo.VoiceConnection, error) {
			{
				return sess.ChannelVoiceJoin(guildID, voiceID, mute, deaf)
			}
		},
		Guild: func(sess *discordgo.Session, inter *discordgo.InteractionCreate) (*discordgo.Guild, error) {
			{
				channel, _ := sess.State.Channel(inter.ChannelID)
				return sess.State.Guild(channel.GuildID)
			}
		},
		Response: func(session *discordgo.Session, interaction *discordgo.InteractionCreate, msg string) {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		},
		Def: &discordgo.ApplicationCommand{
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
					},
					Required: true,
				},
			},
		},
	}

	return &out
}

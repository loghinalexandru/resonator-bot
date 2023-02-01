package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

func NewReact(sync *sync.Map) *types.Playback {
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
							Name:  "Maaaai",
							Value: "misc/mai.dca",
						},
						{
							Name:  "Culcat!",
							Value: "misc/culcat.dca",
						},
						{
							Name:  "Mi-a dat",
							Value: "misc/adat.dca",
						},
						{
							Name:  "Ma bat cainii astia!",
							Value: "misc/caini.dca",
						},
						{
							Name:  "Why are you gay?",
							Value: "misc/whygay.dca",
						},
						{
							Name:  "Noice!",
							Value: "misc/noice.dca",
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

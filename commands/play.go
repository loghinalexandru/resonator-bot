package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

func NewPlay(sync *sync.Map) *types.Playback {
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
			Name:        "play",
			Description: "This command is used to play a sound in the chat!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "type",
					Description: "Sound type to be played!",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Ara-Ara",
							Value: "misc/ara.dca",
						},
						{
							Name:  "Yoooooooouuu",
							Value: "misc/yoo.dca",
						},
						{
							Name:  "FBI Open Up",
							Value: "misc/fbi.dca",
						},
						{
							Name:  "Hehe Boy!",
							Value: "misc/heheboy.dca",
						},
						{
							Name:  "Yamete Kudasai!",
							Value: "misc/yamete.dca",
						},
						{
							Name:  "No God Please No!",
							Value: "misc/nogod.dca",
						},
						{
							Name:  "Mission failed.",
							Value: "misc/failed.dca",
						},
						{
							Name:  "Death",
							Value: "misc/death.dca",
						},
					},
					Required: true,
				},
			},
		},
	}
	return &out
}

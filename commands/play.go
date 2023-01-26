package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/encoder"
)

type playCommand struct {
	identifier string
}

func (cmd playCommand) ID() string {
	return cmd.identifier
}

func (cmd playCommand) Definition() *discordgo.ApplicationCommand {
	result := new(discordgo.ApplicationCommand)
	result.Name = cmd.identifier
	result.Description = "This command is used to play a sound in the chat!"
	result.Options = append(result.Options, &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "type",
		Description: "Sound type to be played!",
		Required:    true,
	})

	return result
}

func (playCommand) Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	channel, _ := session.State.Channel(interaction.ChannelID)
	guild, _ := session.State.Guild(channel.GuildID)

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Playing!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	for _, voice := range guild.VoiceStates {
		if interaction.Member.User.ID == voice.UserID {
			botvc, error := session.ChannelVoiceJoin(guild.ID, voice.ChannelID, false, true)

			if error != nil {
				return error
			}

			path := "misc/" + interaction.ApplicationCommandData().Options[0].Value.(string) + ".mp3"
			stop := make(chan bool)
			encoder.PlayAudioFile(botvc, path, stop)
			break
		}
	}

	return nil
}

package commands

import (
	"github.com/bwmarrin/discordgo"
)

type CustomCommand struct {
	Definition func() *discordgo.ApplicationCommand
	Handler    func(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

func GetCustomCommands() map[string]CustomCommand {
	return map[string]CustomCommand{
		basicCommandID: {
			Definition: describeBasicCommand,
			Handler:    handleBasicCommand,
		},
	}
}

func describeBasicCommand() (result *discordgo.ApplicationCommand) {
	result = new(discordgo.ApplicationCommand)
	result.Name = basicCommandID
	result.Description = "This is a basic command!"

	return
}

func handleBasicCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	channel, _ := session.State.Channel(interaction.ChannelID)
	guild, _ := session.State.Guild(channel.GuildID)

	for _, voice := range guild.VoiceStates {
		if interaction.Member.User.ID == voice.UserID {
			_, error := session.ChannelVoiceJoin(guild.ID, voice.ChannelID, false, true)
			if error != nil {
				return error
			}
		}
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command",
		},
	})

	return nil
}

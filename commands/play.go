package commands

import "github.com/bwmarrin/discordgo"

type playCommand struct {
	identifier string
}

func (cmd playCommand) GetID() string {
	return cmd.identifier
}

func (cmd playCommand) Definition() *discordgo.ApplicationCommand {
	result := new(discordgo.ApplicationCommand)
	result.Name = cmd.identifier
	result.Description = "This command is used to play a sound in the chat!"

	return result
}

func (playCommand) Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
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
			Content: "Pong",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	return nil
}

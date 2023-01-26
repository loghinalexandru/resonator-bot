package commands

import "github.com/bwmarrin/discordgo"

type basicCommand struct {
	identifier string
}

func (cmd basicCommand) ID() string {
	return cmd.identifier
}

func (cmd basicCommand) Definition() (result *discordgo.ApplicationCommand) {
	result = new(discordgo.ApplicationCommand)
	result.Name = cmd.identifier
	result.Description = "Ping!"

	return
}

func (basicCommand) Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong",
		},
	})

	return nil
}

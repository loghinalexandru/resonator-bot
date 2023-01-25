package commands

import "github.com/bwmarrin/discordgo"

var basic = basicCommand{
	identifier: "basic",
}

type basicCommand struct {
	identifier string
}

func (basicCommand) GetID() string {
	return basic.identifier
}

func (basicCommand) Definition() (result *discordgo.ApplicationCommand) {
	result = new(discordgo.ApplicationCommand)
	result.Name = basic.identifier
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

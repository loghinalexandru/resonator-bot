package commands

import "github.com/bwmarrin/discordgo"

var basic = BasicCommand{
	identifier: "basic",
}

type BasicCommand struct {
	identifier string
}

func (BasicCommand) GetID() string {
	return basic.identifier
}

func (BasicCommand) Definition() (result *discordgo.ApplicationCommand) {
	result = new(discordgo.ApplicationCommand)
	result.Name = basic.identifier
	result.Description = "Ping!"

	return
}

func (BasicCommand) Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong",
		},
	})

	return nil
}

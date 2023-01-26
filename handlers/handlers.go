package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands"
)

func Ready(session *discordgo.Session, ready *discordgo.Ready) {
	fmt.Println("Bot is ready!")
	fmt.Println("Bot ShardId: ", session.ShardID)
	fmt.Println("Bot ShardCount: ", session.ShardCount)
}

func InteractionCreate() func(*discordgo.Session, *discordgo.InteractionCreate) {
	commands := commands.Data()

	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		commandError := commands[interaction.ApplicationCommandData().Name].Handler(session, interaction)

		if commandError != nil {
			fmt.Println(commandError)
		}
	}
}

func MessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	//Todo: Add stuff
}

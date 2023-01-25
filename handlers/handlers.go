package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands"
)

func ReadyHandler(session *discordgo.Session, ready *discordgo.Ready) {
	fmt.Println("Bot is ready!")
	fmt.Println("Bot ShardId: ", session.ShardID)
	fmt.Println("Bot ShardCount: ", session.ShardCount)
}

func InteractionCreatedHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	customCommands := commands.GetCustomCommands()

	commandError := customCommands[interaction.ApplicationCommandData().Name].Handler(session, interaction)

	if commandError != nil {
		fmt.Println(commandError)
	}
}

func ReceivedMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	//Todo: Add stuff
}

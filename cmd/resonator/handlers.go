package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Ready() func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(sess *discordgo.Session, inter *discordgo.InteractionCreate) {
		fmt.Println("Bot is ready!")
		fmt.Println("Bot ShardId: ", sess.ShardID)
		fmt.Println("Bot ShardCount: ", sess.ShardCount)
	}
}

func InteractionCreate() func(*discordgo.Session, *discordgo.InteractionCreate) {
	commands := CmdTable()

	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if cmd, ok := commands[interaction.ApplicationCommandData().Name]; ok {
			commandError := cmd.Handler(session, interaction)

			if commandError != nil {
				fmt.Println(commandError)
			}
		}
	}
}

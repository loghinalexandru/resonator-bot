package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func Join(logger *log.Logger) func(*discordgo.Session, *discordgo.GuildCreate) {
	return func(sess *discordgo.Session, gld *discordgo.GuildCreate) {
		logger.Printf("Joined guild with ID %v \n", gld.ID)
	}
}

func InteractionCreate(logger *log.Logger) func(*discordgo.Session, *discordgo.InteractionCreate) {
	commands := CmdTable()

	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if cmd, ok := commands[interaction.ApplicationCommandData().Name]; ok {
			err := cmd.Handler(session, interaction)

			if err != nil {
				logger.Println(err)
			}
		}
	}
}

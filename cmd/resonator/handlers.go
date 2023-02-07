package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Join() func(*discordgo.Session, *discordgo.GuildCreate) {
	return func(sess *discordgo.Session, gld *discordgo.GuildCreate) {
		fmt.Printf("Joined guild with ID %v \n", gld.ID)
	}
}

func InteractionCreate() func(*discordgo.Session, *discordgo.InteractionCreate) {
	commands := CmdTable()

	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if cmd, ok := commands[interaction.ApplicationCommandData().Name]; ok {
			err := cmd.Handler(session, interaction)

			if err != nil {
				//TODO: Add better logging
				fmt.Println(err)
			}
		}
	}
}

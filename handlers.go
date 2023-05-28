package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/logging"
)

func Join(logger *logging.Logger) func(*discordgo.Session, *discordgo.GuildCreate) {
	return func(sess *discordgo.Session, gld *discordgo.GuildCreate) {
		logger.Info("Joined guild with ID: ", gld.ID)
	}
}

func InteractionCreate(cmds []CustomCommandDef, logger *logging.Logger) func(*discordgo.Session, *discordgo.InteractionCreate) {
	var commandsTable = make(map[string]CustomCommandDef)

	for _, cmd := range cmds {
		commandsTable[cmd.Definition().Name] = cmd
	}

	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if cmd, ok := commandsTable[interaction.ApplicationCommandData().Name]; ok {
			err := cmd.Handler(session, interaction)

			if err != nil {
				logger.Error(err)
			}
		}
	}
}

package main

import (
	"github.com/bwmarrin/discordgo"
)

type CustomCommandDef interface {
	Definition() *discordgo.ApplicationCommand
	Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

func Join(logger Logger) func(*discordgo.Session, *discordgo.GuildCreate) {
	return func(sess *discordgo.Session, gld *discordgo.GuildCreate) {
		logger.Info("joined guild", "guildID", gld.ID)
	}
}

func InteractionCreate(cmds []CustomCommandDef, logger Logger) func(*discordgo.Session, *discordgo.InteractionCreate) {
	var commandsTable = make(map[string]CustomCommandDef)

	for _, cmd := range cmds {
		commandsTable[cmd.Definition().Name] = cmd
	}

	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if cmd, ok := commandsTable[interaction.ApplicationCommandData().Name]; ok {
			err := cmd.Handler(session, interaction)

			if err != nil {
				logger.Error("Unexpected application error", "err", err)
			}
		}
	}
}

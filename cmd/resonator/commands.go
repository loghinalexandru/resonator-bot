package main

import (
	"github.com/bwmarrin/discordgo"
)

type CustomCommandDef interface {
	Definition() *discordgo.ApplicationCommand
	Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

func CmdTable() map[string]CustomCommandDef {

	var commandsTable = make(map[string]CustomCommandDef)

	for _, cmd := range cmds {
		commandsTable[cmd.Definition().Name] = cmd
	}

	return commandsTable
}

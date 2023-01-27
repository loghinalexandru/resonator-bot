package commands

import (
	"github.com/bwmarrin/discordgo"
)

var commands = [3]CustomCommandDef{
	playCommand(),
	reactCommand(),
	&basicCommand{
		identifier: "basic",
	},
}

type CustomCommandDef interface {
	Definition() *discordgo.ApplicationCommand
	Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

func Data() map[string]CustomCommandDef {
	var commandsTable = make(map[string]CustomCommandDef)

	for _, cmd := range commands {
		commandsTable[cmd.Definition().Name] = cmd
	}

	return commandsTable
}

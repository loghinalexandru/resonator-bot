package commands

import (
	"github.com/bwmarrin/discordgo"
)

var commands = [2]CustomCommandDef{
	playCommand{
		identifier: "play",
	},
	basicCommand{
		identifier: "basic",
	},
}

type CustomCommandDef interface {
	GetID() string
	Definition() *discordgo.ApplicationCommand
	Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

func GetCustomCommands() map[string]CustomCommandDef {
	var commandsTable = make(map[string]CustomCommandDef)

	for _, cmd := range commands {
		commandsTable[cmd.GetID()] = cmd
	}

	return commandsTable
}

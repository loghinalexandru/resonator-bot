package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

var commands = [2]CustomCommandDef{
	playCommand{
		identifier: "play",
		mutex:      new(sync.Mutex),
	},
	basicCommand{
		identifier: "basic",
	},
}

type CustomCommandDef interface {
	ID() string
	Definition() *discordgo.ApplicationCommand
	Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

func Data() map[string]CustomCommandDef {
	var commandsTable = make(map[string]CustomCommandDef)

	for _, cmd := range commands {
		commandsTable[cmd.ID()] = cmd
	}

	return commandsTable
}

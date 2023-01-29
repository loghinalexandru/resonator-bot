package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	cmdSync  sync.Map
	commands = []CustomCommandDef{
		NewPlay(&cmdSync),
		NewReact(&cmdSync),
		NewRo(&cmdSync),
		NewAnime(),
		NewManga(),
	}
)

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

package main

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/playback"
	"github.com/loghinalexandru/resonator/rest"
)

var (
	cmdSync sync.Map
	cmds    = []CustomCommandDef{
		playback.NewPlay(&cmdSync),
		playback.NewReact(&cmdSync),
		playback.NewRo(&cmdSync),
		rest.NewAnime(),
		rest.NewManga(),
	}
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

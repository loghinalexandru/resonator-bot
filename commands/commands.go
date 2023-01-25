package commands

import (
	"github.com/bwmarrin/discordgo"
)

var customCommands map[string]CustomCommandDef

type CustomCommandDef interface {
	GetID() string
	Definition() *discordgo.ApplicationCommand
	Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

func init() {
	customCommands = map[string]CustomCommandDef{
		basic.identifier: basic,
		play.identifier:  play,
	}
}

func GetCustomCommands() map[string]CustomCommandDef {
	return customCommands
}

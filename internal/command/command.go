package command

import (
	"net/url"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Definition interface {
	Data() *discordgo.ApplicationCommand
	Handle(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

type BotContext struct {
	Sync         *sync.Map
	LogLvl       int
	Token        string
	SwearsApiURL *url.URL
	Index        int
	Shards       int
}

func Join(logger Logger) func(*discordgo.Session, *discordgo.GuildCreate) {
	return func(sess *discordgo.Session, gld *discordgo.GuildCreate) {
		logger.Info("joined guild", "guildID", gld.ID)
	}
}

func InteractionCreate(cmds []Definition, logger Logger) func(*discordgo.Session, *discordgo.InteractionCreate) {
	var commandsTable = make(map[string]Definition)

	for _, cmd := range cmds {
		commandsTable[cmd.Data().Name] = cmd
	}

	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if cmd, ok := commandsTable[interaction.ApplicationCommandData().Name]; ok {
			err := cmd.Handle(session, interaction)

			if err != nil {
				logger.Error("Unexpected application error", "err", err)
			}
		}
	}
}

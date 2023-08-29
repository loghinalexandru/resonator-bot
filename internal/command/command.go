package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/bot"
)

type Handler interface {
	Data() *discordgo.ApplicationCommand
	Handle(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

func Register(sess *discordgo.Session, ctx *bot.Context) error {
	var commandsTable = make(map[string]Handler)

	cmds := []Handler{
		newPlay(ctx),
		newReact(ctx),
		newRo(ctx),
		newCurse(ctx),
		newFeed(ctx),
		newSwear(ctx),
		newAnime(),
		newManga(),
	}

	for _, cmd := range cmds {
		commandsTable[cmd.Data().Name] = cmd
	}

	sess.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if cmd, ok := commandsTable[interaction.ApplicationCommandData().Name]; ok {
			err := cmd.Handle(session, interaction)

			if err != nil {
				ctx.Logger.Error("Unexpected application error", "err", err)
			}
		}
	})

	for _, command := range cmds {
		_, err := sess.ApplicationCommandCreate(sess.State.User.ID, "", command.Data())
		//Remove cmd on termination
		if err != nil {
			ctx.Logger.Error("Unexpected error while creating commands", "err", err)
		}
	}

	return nil
}

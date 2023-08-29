package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/bot"
)

type handler interface {
	Data() *discordgo.ApplicationCommand
	Handle(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

type manager struct {
	commands   []handler
	registered []*discordgo.ApplicationCommand
	ctx        *bot.Context
}

func NewManager(ctx *bot.Context) *manager {
	cmds := []handler{
		newPlay(ctx),
		newReact(ctx),
		newRo(ctx),
		newCurse(ctx),
		newFeed(ctx),
		newSwear(ctx),
		newAnime(),
		newManga(),
	}

	return &manager{
		commands:   cmds,
		registered: make([]*discordgo.ApplicationCommand, len(cmds)),
		ctx:        ctx,
	}
}

func (m *manager) Register(sess *discordgo.Session) {
	var commandsTable = make(map[string]handler)

	for _, cmd := range m.commands {
		commandsTable[cmd.Data().Name] = cmd
	}

	sess.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if cmd, ok := commandsTable[interaction.ApplicationCommandData().Name]; ok {
			err := cmd.Handle(session, interaction)

			if err != nil {
				m.ctx.Logger.Error("Unexpected application error", "err", err)
			}
		}
	})

	for i, cmd := range m.commands {
		r, err := sess.ApplicationCommandCreate(sess.State.User.ID, "", cmd.Data())
		if err != nil {
			m.ctx.Logger.Error("Unexpected error while creating commands", "err", err, "name", cmd.Data().Name)
		}

		m.registered[i] = r
	}
}

func (m *manager) Deregister(sess *discordgo.Session) {
	for _, v := range m.registered {
		err := sess.ApplicationCommandDelete(sess.State.User.ID, "", v.ID)
		if err != nil {
			m.ctx.Logger.Error("Unexpected error while deleting commands", "err", err)
		}
	}
}

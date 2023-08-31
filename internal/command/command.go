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
	commands   map[string]handler
	registered []*discordgo.ApplicationCommand
	ctx        *bot.Context
}

func NewManager(ctx *bot.Context) *manager {
	hh := []handler{
		newPlay(ctx),
		newReact(ctx),
		newRo(ctx),
		newCurse(ctx),
		newFeed(ctx),
		newSwear(ctx),
		newAnime(ctx),
		newManga(ctx),
	}

	cmds := make(map[string]handler, len(hh))

	for _, h := range hh {
		cmds[h.Data().Name] = h
	}

	return &manager{
		commands:   cmds,
		registered: make([]*discordgo.ApplicationCommand, len(hh)),
		ctx:        ctx,
	}
}

func (m *manager) Register(sess *discordgo.Session) {
	sess.AddHandler(m.interactionCreate)

	i := 0
	for _, cmd := range m.commands {
		reg, err := sess.ApplicationCommandCreate(sess.State.User.ID, "", cmd.Data())
		if err != nil {
			m.ctx.Logger.Error("Unexpected error while creating commands", "err", err, "name", cmd.Data().Name)
		}

		m.registered[i] = reg
		i++
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

func (m *manager) interactionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	name := interaction.ApplicationCommandData().Name
	if cmd, ok := m.commands[name]; ok {
		m.ctx.Logger.Info("Handling command", "cmd", name)
		m.ctx.Metrics.ReqCounter.Inc()
		err := cmd.Handle(session, interaction)

		if err != nil {
			m.ctx.Metrics.ErrCounter.Inc()
			m.ctx.Logger.Error("Unexpected application error", "err", err)
		}
	} else {
		m.ctx.Logger.Warn("Could not find handler for command", "name", name)
	}
}

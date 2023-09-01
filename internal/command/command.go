package command

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/bot"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type handler interface {
	Data() *discordgo.ApplicationCommand
	Handle(sess *discordgo.Session, inter *discordgo.InteractionCreate) error
}

type manager struct {
	commands   map[string]handler
	registered []*discordgo.ApplicationCommand
	ctx        *bot.Context
	Request    *prometheus.CounterVec
	Error      *prometheus.CounterVec
	Duration   *prometheus.HistogramVec
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
		Request: promauto.NewCounterVec(prometheus.CounterOpts{
			Subsystem:   "resonator",
			Name:        "command_requests_total",
			Help:        "The total number of commands invoked",
			ConstLabels: prometheus.Labels{"shard": bot.RawID()},
		}, []string{"command"}),
		Error: promauto.NewCounterVec(prometheus.CounterOpts{
			Subsystem:   "resonator",
			Name:        "command_errors_total",
			Help:        "The total number of commands errors",
			ConstLabels: prometheus.Labels{"shard": bot.RawID()},
		}, []string{"command"}),
		Duration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem:   "resonator",
			Name:        "command_duration_seconds",
			Help:        "The duration of a command invocation",
			ConstLabels: prometheus.Labels{"shard": bot.RawID()},
			Buckets:     []float64{.1, .25, .5, 1, 2.5, 5, 10},
		}, []string{"command"}),
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
	cmd, ok := m.commands[name]

	if !ok {
		m.ctx.Logger.Error("Could not find handler for command", "name", name)
		return
	}

	start := time.Now()

	m.Request.With(prometheus.Labels{"command": name}).Inc()
	err := cmd.Handle(session, interaction)

	if err != nil {
		m.ctx.Logger.Error("Unexpected application error", "err", err)
		m.Error.With(prometheus.Labels{"command": name}).Inc()
		return
	}

	m.Duration.With(prometheus.Labels{"command": name}).Observe(time.Since(start).Seconds())
}

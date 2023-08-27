package main

import (
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/command"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func contextFromEnv() command.BotContext {
	ctx := command.BotContext{
		LogLvl: 0,
		Index:  0,
		Shards: 1,
		Sync:   &sync.Map{},
	}

	ctx.Token = os.Getenv("BOT_TOKEN")
	ctx.SwearsApiURL, _ = url.Parse(os.Getenv("SWEARS_API_URL"))

	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		ctx.LogLvl, _ = strconv.Atoi(lvl)
	}

	if replicas := os.Getenv("SHARD_COUNT"); replicas != "" {
		ctx.Shards, _ = strconv.Atoi(replicas)
	}

	if replicaID := os.Getenv("SHARD_ID"); replicaID != "" {
		index := strings.Split(replicaID, "-")
		ctx.Index, _ = strconv.Atoi(index[len(index)-1])
	}

	return ctx
}

func getIntents() discordgo.Intent {
	return discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
}

func main() {
	context := contextFromEnv()

	logLvl := slog.Level(context.LogLvl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLvl,
	}))

	session, err := discordgo.New("Bot " + context.Token)

	if err != nil {
		logger.Error("Unexpected error while creating session", "err", err)
		return
	}

	session.ShouldReconnectVoiceConnOnError = false
	cmds := []command.Definition{
		command.NewPlay(context),
		command.NewReact(context),
		command.NewRo(context),
		command.NewCurse(context),
		command.NewFeed(context),
		command.NewSwear(context),
		command.NewAnime(),
		command.NewManga(),
		command.NewQuote(),
	}

	handlers := []any{
		command.Join(logger),
		command.InteractionCreate(cmds, logger),
	}

	for _, handler := range handlers {
		session.AddHandler(handler)
	}

	session.Identify.Intents = getIntents()
	session.ShardID = context.Index
	session.ShardCount = context.Shards

	err = session.Open()
	defer session.Close()

	if err != nil {
		logger.Error("Unexpected error while opening session", "err", err)
		return
	}

	logger.Info("Bot is ready!")
	logger.Info("Bot shard ID", "shardID", session.ShardID)
	logger.Info("Bot shard count: ", "shardCount", session.ShardCount)

	for _, command := range cmds {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", command.Data())

		if err != nil {
			logger.Error("Unexpected error while creating commands", "err", err)
		}
	}

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigTerm
}

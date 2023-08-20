package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/command"
)

var (
	token        string
	swearsApiURL string
	logLevel     = 0
	shardID      = 0
	shardCount   = 1
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func loadEnv() {
	token = os.Getenv("BOT_TOKEN")
	swearsApiURL = os.Getenv("SWEARS_API_URL")

	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		logLevel, _ = strconv.Atoi(lvl)
	}

	if replicas := os.Getenv("SHARD_COUNT"); replicas != "" {
		shardCount, _ = strconv.Atoi(replicas)
	}

	if replicaID := os.Getenv("SHARD_ID"); replicaID != "" {
		index := strings.Split(replicaID, "-")
		shardID, _ = strconv.Atoi(index[len(index)-1])
	}
}

func getIntents() discordgo.Intent {
	return discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
}

func main() {
	loadEnv()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.Level(logLevel),
	}))

	session, err := discordgo.New("Bot " + token)

	if err != nil {
		logger.Error("Unexpected error while creating session", "err", err)
		return
	}

	session.ShouldReconnectVoiceConnOnError = false
	cmdSync := sync.Map{}
	cmds := []CustomCommandDef{
		command.NewPlay(&cmdSync),
		command.NewReact(&cmdSync),
		command.NewRo(&cmdSync),
		command.NewCurse(&cmdSync, swearsApiURL),
		command.NewFeed(&cmdSync, swearsApiURL),
		command.NewSwear(swearsApiURL),
		command.NewAnime(),
		command.NewManga(),
		command.NewQuote(),
	}

	handlers := []any{
		Join(logger),
		InteractionCreate(cmds, logger),
	}

	for _, handler := range handlers {
		session.AddHandler(handler)
	}

	session.Identify.Intents = getIntents()
	session.ShardID = shardID
	session.ShardCount = shardCount

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
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", command.Definition())

		if err != nil {
			logger.Error("Unexpected error while creating commands", "err", err)
		}
	}

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigTerm
}

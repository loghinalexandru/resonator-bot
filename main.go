package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/bot"
	"github.com/loghinalexandru/resonator/internal/command"
)

func main() {
	ctx := bot.NewContext()
	session, err := discordgo.New("Bot " + bot.Token())

	if err != nil {
		ctx.Logger.Error("Unexpected error while creating session", "err", err)
		return
	}

	session.ShouldReconnectVoiceConnOnError = false
	session.Identify.Intents = bot.Intents()
	session.ShardID = bot.ID()
	session.ShardCount = bot.Shards()

	err = session.Open()
	defer session.Close()

	if err != nil {
		ctx.Logger.Error("Unexpected error while opening session", "err", err)
		return
	}

	bot.Register(session, ctx)
	command.Register(session, ctx)

	ctx.Logger.Info("Bot is ready!")
	ctx.Logger.Info("Bot shard ID", "shardID", session.ShardID)
	ctx.Logger.Info("Bot shard count: ", "shardCount", session.ShardCount)

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigTerm
}

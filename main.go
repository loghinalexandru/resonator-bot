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
	sess, err := discordgo.New("Bot " + bot.Token())

	if err != nil {
		ctx.Logger.Error("Unexpected error while creating session", "err", err)
		return
	}

	sess.ShouldReconnectVoiceConnOnError = false
	sess.Identify.Intents = bot.Intents()
	sess.ShardID = bot.ID()
	sess.ShardCount = bot.Shards()

	err = sess.Open()
	defer sess.Close()

	if err != nil {
		ctx.Logger.Error("Unexpected error while opening session", "err", err)
		return
	}

	cmdManager := command.NewManager(ctx)
	cmdManager.Register(sess)

	if err != nil {
		panic(err)
	}

	ctx.Logger.Info("Bot is ready!", "shardID", sess.ShardID, "shardCount", sess.ShardCount)

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigTerm

	ctx.Logger.Info("Bot is shutting down!")

	if bot.Cleanup() {
		ctx.Logger.Info("Deregistering commands!")
		cmdManager.Deregister(sess)
	}
}

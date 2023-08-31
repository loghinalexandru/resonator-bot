package main

import (
	"context"
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

	ctx.Logger.Info("Bot is ready", "shardID", sess.ShardID, "shardCount", sess.ShardCount)

	shutdown := bot.StartMetricsServer(ctx)
	defer shutdown()

	s, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s.Done()

	ctx.Logger.Info("Bot is shutting down")

	if bot.Cleanup() {
		ctx.Logger.Info("Deregistering commands")
		cmdManager.Deregister(sess)
	}
}

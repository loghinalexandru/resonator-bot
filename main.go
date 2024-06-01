package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/bot"
	"github.com/loghinalexandru/resonator/internal/command"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func main() {
	ctx := bot.NewContext()
	cmdManager := command.NewManager(ctx)
	registry := prometheus.NewRegistry()

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		cmdManager.Error,
		cmdManager.Request,
		cmdManager.Duration,
	)

	sess, err := discordgo.New("Bot " + bot.Token())
	if err != nil {
		ctx.Logger.Error("Unexpected error while creating session", "err", err)
		return
	}

	sess.ShouldReconnectVoiceOnSessionError = false
	sess.Identify.Intents = bot.Intents()
	sess.ShardID = bot.ID()
	sess.ShardCount = bot.Shards()

	err = sess.Open()
	if err != nil {
		ctx.Logger.Error("Unexpected error while opening session", "err", err)
		return
	}

	defer sess.Close()

	shutdown := bot.StartMetricsServer(ctx.Logger, registry)
	defer shutdown()

	cmdManager.Register(sess)
	ctx.Logger.Info("Bot is ready", "shardID", sess.ShardID, "shardCount", sess.ShardCount)

	s, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s.Done()

	ctx.Logger.Info("Bot is shutting down")

	if bot.Cleanup() {
		ctx.Logger.Info("Unregistering commands")
		cmdManager.Deregister(sess)
	}
}

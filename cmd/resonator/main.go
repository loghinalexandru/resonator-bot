package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/commands"
	"github.com/loghinalexandru/resonator/pkg/logging"
)

var (
	token        string
	swearsApiURL string
	logLevel     logging.LogLevel
	shardID      int = 0
	shardCount   int = 1
)

func loadEnv() {
	token = os.Getenv("BOT_TOKEN")
	swearsApiURL = os.Getenv("SWEARS_API_URL")
	logLevel = logging.ToLogLevel(os.Getenv("LOG_LEVEL"))

	shardCount, _ = strconv.Atoi(os.Getenv("SHARD_COUNT"))
	index := strings.Split(os.Getenv("SHARD_ID"), "-")
	shardID, _ = strconv.Atoi(index[len(index)-1])
}

func getIntents() discordgo.Intent {
	return discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
}

func main() {
	loadEnv()

	session, sessionError := discordgo.New("Bot " + token)
	session.ShouldReconnectVoiceConnOnError = false

	logger := logging.New(logLevel, log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile))

	cmdSync := sync.Map{}
	cmds := []CustomCommandDef{
		commands.NewPlay(&cmdSync),
		commands.NewReact(&cmdSync),
		commands.NewRo(&cmdSync),
		commands.NewCurse(&cmdSync, swearsApiURL),
		commands.NewSwear(swearsApiURL, http.DefaultClient),
		commands.NewAnime(http.DefaultClient),
		commands.NewManga(http.DefaultClient),
		commands.NewQuote(http.DefaultClient),
	}

	handlers := []any{
		Join(logger),
		InteractionCreate(cmds, logger),
	}

	if sessionError != nil {
		logger.Error(sessionError)
		return
	}

	for _, handler := range handlers {
		session.AddHandler(handler)
	}

	session.Identify.Intents = getIntents()
	session.ShardID = shardID
	session.ShardCount = shardCount

	socketError := session.Open()
	defer session.Close()

	if socketError != nil {
		logger.Error(socketError)
		return
	}

	logger.Info("Bot is ready!")
	logger.Info("Bot ShardId: ", session.ShardID)
	logger.Info("Bot ShardCount: ", session.ShardCount)

	for _, command := range cmds {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", command.Definition())

		if err != nil {
			logger.Error(err)
		}
	}

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigTerm
}

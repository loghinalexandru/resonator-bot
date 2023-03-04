package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	lvl, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	logLevel = logging.LogLevel(lvl)
}

func getIntents() discordgo.Intent {
	return discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
}

func main() {
	loadEnv()

	session, sessionError := discordgo.New("Bot " + token)
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

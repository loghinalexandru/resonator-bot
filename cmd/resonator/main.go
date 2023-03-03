package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/commands"
	"github.com/loghinalexandru/resonator/pkg/logging"
)

func readEnv() (string, string) {
	token, success := os.LookupEnv("BOT_TOKEN")

	if !success {
		token = ""
	}

	swearsApiURL, success := os.LookupEnv("SWEARS_API_URL")

	if !success {
		swearsApiURL = ""
	}

	return token, swearsApiURL
}

func getIntents() discordgo.Intent {
	return discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
}

func main() {
	token, swearsApiURL := readEnv()

	session, sessionError := discordgo.New("Bot " + token)
	logger := logging.New(logging.Info, log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile))

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

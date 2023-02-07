package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/commands"
)

var (
	logger       *log.Logger
	token        string
	swearsApiURL string
	cmdSync      sync.Map
	cmds         []CustomCommandDef
)

func init() {
	logger = log.Default()
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	var success bool
	token, success = os.LookupEnv("BOT_TOKEN")

	if !success {
		token = ""
	}

	swearsApiURL, success = os.LookupEnv("SWEARS_API_URL")

	if !success {
		swearsApiURL = ""
	}

	cmds = []CustomCommandDef{
		commands.NewPlay(&cmdSync),
		commands.NewReact(&cmdSync),
		commands.NewRo(&cmdSync),
		commands.NewCurse(&cmdSync, swearsApiURL),
		commands.NewSwear(swearsApiURL),
		commands.NewAnime(),
		commands.NewManga(),
		commands.NewQuote(),
	}
}

func getIntents() discordgo.Intent {
	return discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
}

func getHandlers() []interface{} {
	return []any{
		Join(logger),
		InteractionCreate(logger),
	}
}

func main() {
	session, sessionError := discordgo.New("Bot " + token)

	if sessionError != nil {
		logger.Println(sessionError)
		return
	}

	for _, handler := range getHandlers() {
		session.AddHandler(handler)
	}

	session.Identify.Intents = getIntents()

	socketError := session.Open()
	defer session.Close()

	if socketError != nil {
		logger.Println(socketError)
		return
	}

	logger.Println("Bot is ready!")
	logger.Println("Bot ShardId: ", session.ShardID)
	logger.Println("Bot ShardCount: ", session.ShardCount)

	for _, command := range CmdTable() {
		_, err := session.ApplicationCommandCreate(
			session.State.User.ID, "", command.Definition())

		if err != nil {
			logger.Println(err)
		}
	}

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigTerm
}

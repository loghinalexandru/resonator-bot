package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/commands"
)

var (
	token     string
	swearsURL string
	cmdSync   sync.Map
	cmds      []CustomCommandDef
)

func init() {
	var success bool
	token, success = os.LookupEnv("BOT_TOKEN")

	if !success {
		token = ""
	}

	swearsURL, success = os.LookupEnv("SWEARS_API")

	if !success {
		swearsURL = ""
	}

	cmds = []CustomCommandDef{
		commands.NewPlay(&cmdSync),
		commands.NewReact(&cmdSync),
		commands.NewRo(&cmdSync),
		commands.NewCurse(&cmdSync, swearsURL),
		commands.NewSwear(swearsURL),
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
		Ready(),
		InteractionCreate(),
	}
}

func main() {
	session, sessionError := discordgo.New("Bot " + token)

	if sessionError != nil {
		fmt.Println(sessionError)
		return
	}

	for _, handler := range getHandlers() {
		session.AddHandler(handler)
	}

	session.Identify.Intents = getIntents()

	socketError := session.Open()
	defer session.Close()

	if socketError != nil {
		fmt.Println(socketError)
		return
	}

	for _, command := range CmdTable() {
		_, err := session.ApplicationCommandCreate(
			session.State.User.ID, "", command.Definition())

		if err != nil {
			fmt.Println(err)
		}
	}

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigTerm
}

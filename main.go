package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands"
	"github.com/loghinalexandru/resonator/handlers"
)

var token string

func init() {
	var success bool
	token, success = os.LookupEnv("BOT_TOKEN")

	if !success {
		token = ""
	}
}

func getIntents() discordgo.Intent {
	return discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
}

func getHandlers() []interface{} {
	return []any{
		handlers.ReceivedMessageHandler,
		handlers.ReadyHandler,
		handlers.InteractionCreatedHandler,
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

	for _, command := range commands.GetCustomCommands() {
		_, commandCreationError := session.ApplicationCommandCreate(session.State.User.ID, "", command.Definition())

		if commandCreationError != nil {
			fmt.Println(commandCreationError)
		}
	}

	signalTermination := make(chan os.Signal)
	signal.Notify(signalTermination, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-signalTermination
}

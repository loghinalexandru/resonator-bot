package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var token string = ""

func receivedMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	fmt.Println(message.ChannelID)
	fmt.Println(message.Application)
	fmt.Println(message.Author)
	fmt.Println(message.Content)
}

func main() {
	bot, _ := discordgo.New("Bot " + token)
	bot.Identify.Intents = discordgo.IntentsGuildMessages

	bot.AddHandler(receivedMessage)
	error := bot.Open()

	if error != nil {
		fmt.Println(error)
	}

	fmt.Println("Bot Started!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

package command

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/rest"
)

type quoteData struct {
	Quote string `json:"quote"`
}

func NewQuote() *rest.REST[quoteData] {
	url := "https://api.kanye.rest/"
	def := &discordgo.ApplicationCommand{
		Name:        "quote",
		Description: "This command is used find Kanye West quotes!",
	}

	result, err := rest.New(def, url, rest.WithFormatter[quoteData](quoteFormatter))

	if err != nil {
		panic(err)
	}

	return result
}

func quoteFormatter(content quoteData) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("> \"**%s**\" - Kanye West", content.Quote))
	return sb.String()
}

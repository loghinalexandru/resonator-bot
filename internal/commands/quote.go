package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/rest"
)

type quoteData struct {
	Quote string `json:"quote"`
}

func NewQuote() *rest.REST {
	out := rest.New(&discordgo.ApplicationCommand{
		Name:        "quote",
		Description: "This command is used find Kanye West quotes!",
	},
		rest.WithURL("https://api.kanye.rest/"),
		rest.WithDataType(&quoteData{}),
		rest.WithFormatter(quoteFormatter),
	)

	return &out
}

func quoteFormatter(content any) string {
	var sb strings.Builder
	resp, ok := content.(*quoteData)

	if !ok {
		return "Something went wrong!"
	}

	sb.WriteString(fmt.Sprintf("> \"**%s**\" - Kanye West", resp.Quote))
	return sb.String()
}

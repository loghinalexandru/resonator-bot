package commands

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/rest"
)

type quoteData struct {
	Quote string `json:"quote"`
}

func NewQuote(client *http.Client) *rest.REST[quoteData] {
	url := "https://api.kanye.rest/"
	def := &discordgo.ApplicationCommand{
		Name:        "quote",
		Description: "This command is used find Kanye West quotes!",
	}

	return rest.New(def, url, client, quoteFormatter)
}

func quoteFormatter(content quoteData) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("> \"**%s**\" - Kanye West", content.Quote))
	return sb.String()
}

package rest

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type quoteData struct {
	Quote string `json:"quote"`
}

func NewQuote() *REST {
	out := REST{
		URL:       "https://api.kanye.rest/",
		Type:      &quoteData{},
		Formatter: quoteFormatter,
		Def: &discordgo.ApplicationCommand{
			Name:        "quote",
			Description: "This command is used find Kanye West quotes!",
		},
	}
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

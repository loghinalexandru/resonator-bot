package command

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/bot"
	"github.com/loghinalexandru/resonator/pkg/rest"
)

const (
	animeURL = "https://kitsu.io/api/edge/anime?filter[text]=%v&page[limit]=10"
)

type animeData struct {
	ID    string `json:"id"`
	Stats struct {
		Title    string `json:"canonicalTitle"`
		Type     string `json:"showType"`
		Rated    string `json:"ageRating"`
		Status   string `json:"status"`
		Start    string `json:"startDate"`
		End      string `json:"endDate"`
		Episodes int32  `json:"episodeCount"`
		Length   int32  `json:"episodeLength"`
		Total    int32  `json:"totalLength"`
	} `json:"attributes"`
}

type animeWrapper struct {
	Content []animeData `json:"data"`
}

func newAnime(ctx *bot.Context) *rest.REST[animeWrapper] {
	def := &discordgo.ApplicationCommand{
		Name:        "anime",
		Description: "This command is used find anime via Kitsu API!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "keyword",
				Description: "Keyword to search for: ",
				Required:    true,
			},
		},
	}

	result, err := rest.New(def, animeURL, rest.WithFormatter[animeWrapper](animeFormatter))

	if err != nil {
		ctx.Logger.Error("Error creating /anime command", "err", err)
	}

	return result
}

func animeFormatter(resp animeWrapper) string {
	var sb strings.Builder

	if len(resp.Content) > 0 {
		sb.WriteString("Here's a list: \n")
		for _, r := range resp.Content {
			sb.WriteString(fmt.Sprintf("**%s - %s**\n", r.Stats.Title, strings.ToUpper(r.Stats.Status)))
			sb.WriteString(fmt.Sprintf("> Type: %v, Episodes: %v, Length:  %vm, Total Time: %vm\n", r.Stats.Type, r.Stats.Episodes, r.Stats.Length, r.Stats.Total))
		}
	} else {
		sb.WriteString("No match found :(")
	}

	return sb.String()
}

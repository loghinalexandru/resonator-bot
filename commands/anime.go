package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/commands/types"
)

type content struct {
	Id    string `json:"id"`
	Stats struct {
		Title    string `json:"canonicalTitle"`
		Type     string `json:"showType"`
		Rating   string `json:"averageRating"`
		Rated    string `json:"ageRating"`
		Status   string `json:"status"`
		Start    string `json:"startDate"`
		End      string `json:"endDate"`
		Episodes int32  `json:"episodeCount"`
		Length   int32  `json:"episodeLength"`
		Total    int32  `json:"totalLength"`
	} `json:"attributes"`
}

type wrapper struct {
	Content []content `json:"data"`
}

func NewAnime() *types.REST {
	out := types.REST{
		URL:       "https://kitsu.io/api/edge/anime?filter[text]=%v&page[limit]=5",
		Type:      &wrapper{},
		Formatter: kitsuFormat,
		Def: &discordgo.ApplicationCommand{
			Name:        "anime",
			Description: "This command is used find anime via kitsu API!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keyword",
					Description: "Keyword to search for",
					Required:    true,
				},
			},
		},
	}
	return &out
}

func kitsuFormat(content any) string {
	var sb strings.Builder
	resp, ok := content.(*wrapper)

	if !ok {
		return "Something went wrong!"
	}

	if len(resp.Content) > 0 {
		sb.WriteString("Here's a list: \n")
		for _, r := range resp.Content {
			sb.WriteString(fmt.Sprintf("**%s - %s**\n", r.Stats.Title, strings.ToUpper(r.Stats.Status)))
			sb.WriteString(fmt.Sprintf("> Episodes: %v, Length:  %vm, Total Time: %vm\n", r.Stats.Episodes, r.Stats.Length, r.Stats.Total))
		}
	} else {
		sb.WriteString("No match found :(")
	}

	return sb.String()
}

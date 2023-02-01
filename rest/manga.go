package rest

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type mangaData struct {
	Id    string `json:"id"`
	Stats struct {
		Title       string `json:"canonicalTitle"`
		Type        string `json:"showType"`
		Rated       string `json:"ageRating"`
		RatingGuide string `json:"ageRatingGuide"`
		Status      string `json:"status"`
		Start       string `json:"startDate"`
		End         string `json:"endDate"`
		Chapters    int32  `json:"chapterCount"`
		Volumes     int32  `json:"volumeCount"`
	} `json:"attributes"`
}

type mangaWrapper struct {
	Content []mangaData `json:"data"`
}

func NewManga() *REST {
	out := REST{
		URL:       "https://kitsu.io/api/edge/manga?filter[text]=%v&filter[subtype]=manga&page[limit]=5",
		Type:      &mangaWrapper{},
		Formatter: mangaFormatter,
		Def: &discordgo.ApplicationCommand{
			Name:        "manga",
			Description: "This command is used find manga via Kitsu API!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keyword",
					Description: "Keyword to search for: ",
					Required:    true,
				},
			},
		},
	}
	return &out
}

func mangaFormatter(content any) string {
	var sb strings.Builder
	resp, ok := content.(*mangaWrapper)

	if !ok {
		return "Something went wrong!"
	}

	if len(resp.Content) > 0 {
		sb.WriteString("Here's a list: \n")
		for _, r := range resp.Content {
			sb.WriteString(fmt.Sprintf("**%s - %s**\n", r.Stats.Title, strings.ToUpper(r.Stats.Status)))
			sb.WriteString(fmt.Sprintf("> Volumes: %v, Chapters: %v, Genre: %v\n", r.Stats.Volumes, r.Stats.Chapters, r.Stats.RatingGuide))
		}
	} else {
		sb.WriteString("No match found :(")
	}

	return sb.String()
}

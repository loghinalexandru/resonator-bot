package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Kitsu struct {
	URL string
	Def *discordgo.ApplicationCommand
}

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

func (cmd *Kitsu) Definition() *discordgo.ApplicationCommand {
	return cmd.Def
}

func (cmd *Kitsu) Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error {
	customURL := fmt.Sprintf(cmd.URL, inter.ApplicationCommandData().Options[0].Value)
	response, err := http.Get(customURL)

	if err != nil {
		return nil
	}

	defer response.Body.Close()
	content, _ := io.ReadAll(response.Body)

	var resp wrapper
	var sb strings.Builder
	err = json.Unmarshal(content, &resp)

	if err != nil {
		fmt.Println(err)
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

	sess.InteractionRespond(inter.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
		},
	})

	return nil
}

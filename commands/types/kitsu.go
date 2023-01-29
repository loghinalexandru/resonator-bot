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

type data struct {
	Id        string `json:"id"`
	EntryType string `json:"type"`
	Stats     struct {
		Title  string `json:"canonicalTitle"`
		Rating string `json:"averageRating"`
		Rated  string `json:"ageRating"`
		Status string `json:"status"`
		Start  string `json:"startDate"`
		End    string `json:"endDate"`
	} `json:"attributes"`
}

type kitsuResp struct {
	Content []data `json:"data"`
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

	var resp kitsuResp
	var sb strings.Builder
	err = json.Unmarshal(content, &resp)

	if err != nil {
		fmt.Println(err)
	}
	if len(resp.Content) > 0 {
		sb.WriteString("Here's a list: \n")
		for _, r := range resp.Content {
			sb.WriteString(fmt.Sprintf("**%v** %v, First Aired: %v, Ended In: %v", r.Stats.Title, r.Stats.Status, r.Stats.Start, r.Stats.End))
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("No match found :(")
	}

	sess.InteractionRespond(inter.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	return nil
}

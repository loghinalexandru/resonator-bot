package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var (
	respond = sendResp
)

const (
	failure = "Service can not be reached!"
)

type REST[T any] struct {
	url       string
	data      T
	formatter func(payload T) string
	client    *http.Client
	def       *discordgo.ApplicationCommand
}

func New[T any](definition *discordgo.ApplicationCommand, url string, client *http.Client, form func(payload T) string) *REST[T] {
	result := REST[T]{
		def:       definition,
		url:       url,
		client:    client,
		formatter: form,
	}

	return &result
}

func (cmd *REST[T]) Definition() *discordgo.ApplicationCommand {
	return cmd.def
}

func (cmd *REST[T]) Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var args []any
	for _, v := range inter.ApplicationCommandData().Options {
		if v.Type == discordgo.ApplicationCommandOptionString {
			args = append(args, v.Value)
		}
	}

	response, err := cmd.client.Get(fmt.Sprintf(cmd.url, args...))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		sess.InteractionRespond(inter.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: failure,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		return nil
	}

	defer response.Body.Close()
	content, _ := io.ReadAll(response.Body)

	err = json.Unmarshal(content, &cmd.data)

	if err != nil {
		return err
	}

	respond(sess, inter.Interaction, cmd.createReponse())

	return nil
}

// Seam functions for testing
func (cmd *REST[T]) createReponse() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: cmd.formatter(cmd.data),
		},
	}
}

func sendResp(sess *discordgo.Session, inter *discordgo.Interaction, resp *discordgo.InteractionResponse) {
	sess.InteractionRespond(inter, resp)
}

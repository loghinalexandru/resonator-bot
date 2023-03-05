package rest

import (
	"encoding/json"
	"errors"
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
	var data T
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
		interResp := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: failure,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		}

		respond(sess, inter.Interaction, interResp)
		return errors.New("call to URI failed")
	}

	defer response.Body.Close()
	content, _ := io.ReadAll(response.Body)

	err = json.Unmarshal(content, &data)

	if err != nil {
		return err
	}

	respond(sess, inter.Interaction, cmd.createReponse(data))
	return nil
}

// Seam functions for testing
func (cmd *REST[T]) createReponse(data T) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: cmd.formatter(data),
		},
	}
}

func sendResp(sess *discordgo.Session, inter *discordgo.Interaction, resp *discordgo.InteractionResponse) {
	sess.InteractionRespond(inter, resp)
}

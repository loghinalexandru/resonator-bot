package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

const (
	failure = "Service can not be reached!"
)

type REST[T any] struct {
	url       string
	data      T
	formatter func(payload T) string
	def       *discordgo.ApplicationCommand
}

func New[T any](definition *discordgo.ApplicationCommand, url string, form func(payload T) string) *REST[T] {
	result := REST[T]{
		def:       definition,
		url:       url,
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

	response, err := http.Get(fmt.Sprintf(cmd.url, args...))
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

	sess.InteractionRespond(inter.Interaction, cmd.createReponse())
	return nil
}

func (cmd *REST[T]) createReponse() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: cmd.formatter(cmd.data),
		},
	}
}

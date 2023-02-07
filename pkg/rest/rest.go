package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type OptionsFunc func(cmd *REST)

const (
	failure = "Service can not be reached!"
)

type REST struct {
	url       string
	dataType  any
	formatter func(payload any) string
	def       *discordgo.ApplicationCommand
}

func New(definition *discordgo.ApplicationCommand, options ...OptionsFunc) REST {
	result := REST{
		def: definition,
	}

	for _, opt := range options {
		opt(&result)
	}

	return result
}

func WithURL(url string) OptionsFunc {
	return func(cmd *REST) {
		cmd.url = url
	}
}

func WithDataType(dataType any) OptionsFunc {
	return func(cmd *REST) {
		cmd.dataType = dataType
	}
}

func WithFormatter(formatter func(payload any) string) OptionsFunc {
	return func(cmd *REST) {
		cmd.formatter = formatter
	}
}

func (cmd *REST) Definition() *discordgo.ApplicationCommand {
	return cmd.def
}

func (cmd *REST) Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var args []any
	for _, v := range inter.ApplicationCommandData().Options {
		args = append(args, v.Value.(string))
	}

	customURL := fmt.Sprintf(cmd.url, args...)
	response, err := http.Get(customURL)

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

	err = json.Unmarshal(content, cmd.dataType)

	if err != nil {
		return err
	}

	sess.InteractionRespond(inter.Interaction, cmd.createReponse())
	return nil
}

func (cmd *REST) createReponse() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: cmd.formatter(cmd.dataType),
		},
	}
}

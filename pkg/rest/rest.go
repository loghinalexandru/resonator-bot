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
	ErrCallFailed    = errors.New("call to URI failed")
	ErrHttpClient    = errors.New("missing http client")
	ErrRespFormatter = errors.New("missing response formatter")
	respond          = discordgoResp
)

const (
	msgCallFailed = "Service can not be reached!"
)

type restOpt[T any] func(*REST[T]) error
type respFmt[T any] func(payload T) string

type REST[T any] struct {
	baseURL   string
	formatter respFmt[T]
	client    *http.Client
	def       *discordgo.ApplicationCommand
}

func New[T any](definition *discordgo.ApplicationCommand, URL string, opts ...restOpt[T]) (*REST[T], error) {
	result := &REST[T]{
		def:       definition,
		baseURL:   URL,
		client:    http.DefaultClient,
		formatter: func(p T) string { return fmt.Sprint(p) },
	}

	for _, opt := range opts {
		err := opt(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func WithHttpClient[T any](c *http.Client) restOpt[T] {
	return func(r *REST[T]) error {
		if r.client == nil {
			return ErrHttpClient
		}

		r.client = c
		return nil
	}
}

func WithFormatter[T any](f respFmt[T]) restOpt[T] {
	return func(r *REST[T]) error {
		if r.client == nil {
			return ErrRespFormatter
		}

		r.formatter = f
		return nil
	}
}

func (cmd *REST[T]) Data() *discordgo.ApplicationCommand {
	return cmd.def
}

func (cmd *REST[T]) Handle(sess *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var args []any
	var data T
	for _, v := range inter.ApplicationCommandData().Options {
		if v.Type == discordgo.ApplicationCommandOptionString {
			args = append(args, v.Value)
		}
	}

	response, err := cmd.client.Get(fmt.Sprintf(cmd.baseURL, args...))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		r := createReponse(msgCallFailed, discordgo.MessageFlagsEphemeral)
		respond(sess, inter.Interaction, r)
		return ErrCallFailed
	}

	defer response.Body.Close()
	content, _ := io.ReadAll(response.Body)

	err = json.Unmarshal(content, &data)

	if err != nil {
		return err
	}

	respond(sess, inter.Interaction, createReponse(cmd.formatter(data)))
	return nil
}

func createReponse(data string, flags ...discordgo.MessageFlags) *discordgo.InteractionResponse {
	var result discordgo.MessageFlags

	for _, f := range flags {
		result = result | f
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: data,
			Flags:   result,
		},
	}
}

// Seam functions for testing
func discordgoResp(sess *discordgo.Session, inter *discordgo.Interaction, resp *discordgo.InteractionResponse) {
	sess.InteractionRespond(inter, resp)
}

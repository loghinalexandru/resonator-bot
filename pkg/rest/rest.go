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
	ErrRespFormatter = errors.New("missing response formatter")
	respond          = discordgoResp
)

const (
	msgCallFailed = "Service can not be reached!"
)

type interactionResp func(sess *discordgo.Session, inter *discordgo.Interaction, resp *discordgo.InteractionResponse) error

type Opt[T any] func(*REST[T]) error
type Formatter[T any] func(payload T) string

type REST[T any] struct {
	baseURL   string
	formatter Formatter[T]
	client    *http.Client
	def       *discordgo.ApplicationCommand
	resp      interactionResp
}

func New[T any](definition *discordgo.ApplicationCommand, url string, opts ...Opt[T]) (*REST[T], error) {
	return newInternal(definition, url, discordgoResp, opts...)
}

func newInternal[T any](definition *discordgo.ApplicationCommand, url string, resp interactionResp, opts ...Opt[T]) (*REST[T], error) {
	result := &REST[T]{
		def:       definition,
		baseURL:   url,
		client:    http.DefaultClient,
		formatter: func(p T) string { return fmt.Sprint(p) },
		resp:      resp,
	}

	for _, opt := range opts {
		err := opt(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func WithFormatter[T any](f Formatter[T]) Opt[T] {
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
		r := createResponse(msgCallFailed, discordgo.MessageFlagsEphemeral)
		err = cmd.resp(sess, inter.Interaction, r)

		return errors.Join(err, ErrCallFailed)
	}

	defer response.Body.Close()
	content, _ := io.ReadAll(response.Body)

	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}

	err = cmd.resp(sess, inter.Interaction, createResponse(cmd.formatter(data)))
	if err != nil {
		return err
	}

	return nil
}

func createResponse(data string, flags ...discordgo.MessageFlags) *discordgo.InteractionResponse {
	var result discordgo.MessageFlags

	for _, f := range flags {
		result |= f
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: data,
			Flags:   result,
		},
	}
}

func discordgoResp(sess *discordgo.Session, inter *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	return sess.InteractionRespond(inter, resp)
}

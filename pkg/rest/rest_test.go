package rest

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/bwmarrin/discordgo"
)

const (
	tstMessage = "test"
)

type TestStruct struct {
}

func TestHandler(t *testing.T) {
	t.Parallel()

	respond = sendRespMock
	var testPayload TestStruct
	var cmdInter *discordgo.InteractionCreate

	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("{}")),
			Header:     make(http.Header),
		}
	})

	cmdInter = &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommandAutocomplete,
			Data: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{},
			},
		},
	}

	target := REST[TestStruct]{
		data:   testPayload,
		client: client,
		formatter: func(payload TestStruct) string {
			return tstMessage
		},
	}

	err := target.Handler(&discordgo.Session{}, cmdInter)

	if err != nil {
		t.Error("Should not be nil!")
	}
}

func TestCreateResponse(t *testing.T) {
	t.Parallel()

	var testPayload TestStruct
	target := &REST[TestStruct]{
		data: testPayload,
		formatter: func(payload TestStruct) string {
			return tstMessage
		},
	}

	res := target.createReponse()

	if res.Data.Content != tstMessage {
		t.Error("Reponse message does not match!")
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func sendRespMock(sess *discordgo.Session, inter *discordgo.Interaction, resp *discordgo.InteractionResponse) {
	return
}

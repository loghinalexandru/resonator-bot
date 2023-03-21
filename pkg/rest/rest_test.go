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

func TestNew(t *testing.T) {
	testDef := &discordgo.ApplicationCommand{
		ID: "testDef",
	}
	testURL := "test"

	got := New(
		testDef,
		testURL,
		http.DefaultClient,
		func(payload TestStruct) string {
			return "test"
		})

	if got.def != testDef {
		t.Error("different command definition")
	}

	if got.url != testURL {
		t.Error("different command url")
	}

	if got.client != http.DefaultClient {
		t.Error("different command http client")
	}

	if got.formatter == nil {
		t.Error("different command formatter")
	}
}

func TestDefinition(t *testing.T) {
	t.Parallel()

	target := &REST[TestStruct]{
		def: &discordgo.ApplicationCommand{},
	}

	if target.Definition() == nil {
		t.Error("this should not be nil")
	}
}

func TestHandlerWhenCallFails(t *testing.T) {
	t.Parallel()

	respond = sendRespMock
	var cmdInter *discordgo.InteractionCreate

	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 500,
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
		client: client,
		formatter: func(payload TestStruct) string {
			return tstMessage
		},
	}

	err := target.Handler(&discordgo.Session{}, cmdInter)

	if err == nil || err.Error() != "call to URI failed" {
		t.Error(err)
	}
}

func TestHandler(t *testing.T) {
	t.Parallel()

	respond = sendRespMock
	var cmdInter *discordgo.InteractionCreate

	client := newTestClient(func(req *http.Request) *http.Response {
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
		client: client,
		formatter: func(payload TestStruct) string {
			return tstMessage
		},
	}

	err := target.Handler(&discordgo.Session{}, cmdInter)

	if err != nil {
		t.Error(err)
	}
}

func TestCreateResponse(t *testing.T) {
	t.Parallel()

	var testPayload TestStruct
	target := &REST[TestStruct]{
		formatter: func(payload TestStruct) string {
			return tstMessage
		},
	}

	got := target.createReponse(testPayload)

	if got.Data.Content != tstMessage {
		t.Errorf("want %q, got %q", tstMessage, got.Data.Content)
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func sendRespMock(sess *discordgo.Session, inter *discordgo.Interaction, resp *discordgo.InteractionResponse) {
}

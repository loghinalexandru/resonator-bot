package playback

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestNew(t *testing.T) {
	t.Parallel()

	testDef := &discordgo.ApplicationCommand{
		ID: "testDef",
	}

	testSyncMap := &sync.Map{}

	got := New(testSyncMap, testDef)

	if got.def != testDef {
		t.Error("different command definition")
	}

	if got.storage != testSyncMap {
		t.Error("different command storage")
	}
}

func TestDefinition(t *testing.T) {
	t.Parallel()

	target := &Playback{
		def: &discordgo.ApplicationCommand{},
	}

	if target.Definition() == nil {
		t.Error("this should not be nil")
	}
}

func TestPlaySoundWithValidFile(t *testing.T) {
	t.Parallel()

	testChan := make(chan []byte, 100)
	fh, _ := os.Open("testdata/test_file.dca")
	defer fh.Close()

	err := playSound(testChan, fh)

	if err != nil {
		t.Fatal(err)
	}

	packet := <-testChan

	if packet == nil {
		t.Fatal("this should not be nil")
	}
}

func TestPlaySoundTimeouts(t *testing.T) {
	t.Parallel()

	testChan := make(chan []byte)
	fh, _ := os.Open("testdata/test_file.dca")
	defer fh.Close()

	err := playSound(testChan, fh)

	if err == nil {
		t.Fatal("no timeout error")
	}
}

func TestPlaySoundWithNilHandler(t *testing.T) {
	t.Parallel()

	err := playSound(make(chan<- []byte), nil)

	if err == nil {
		t.Error("this should not be nil")
	}
}

func TestPlaySoundWithInvalidHandler(t *testing.T) {
	t.Parallel()

	fh, _ := os.Open("")
	err := playSound(make(chan<- []byte), fh)

	if err == nil {
		t.Error("this should not be nil")
	}
}

func TestGetAudioSourceWithInvalidURI(t *testing.T) {
	t.Parallel()

	tests := []string{"", "localhost.com/api", "//google", "http://"}

	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			_, err := getAudioSource(tc, nil)

			if err == nil {
				t.Error(err)
			}
		})
	}
}

func TestGetAudioSourceWithValidURI(t *testing.T) {
	t.Parallel()

	tests := []string{"testdata/test_file.dca", "https://www.google.com/"}

	client := newTestClient(func(req *http.Request) *http.Response {
		if req.URL.String() == "https://www.google.com/" {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("{}")),
				Header:     make(http.Header),
			}
		}
		return nil
	})

	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			res, err := getAudioSource(tc, client)

			if err != nil {
				t.Fatal(err)
			}

			if res == nil {
				t.Error("should not be nil")
			}
		})
	}
}

func TestHandlerWhenCalledCreatesMutex(t *testing.T) {
	t.Parallel()

	voice = joinVoiceMock
	guild = getGuildMock
	respond = sendRespMock

	inter := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Member: &discordgo.Member{
				User: &discordgo.User{
					ID: "user",
				},
			},
			Data: discordgo.ApplicationCommandInteractionData{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Value: "test",
					},
				},
			},
		},
	}
	target := &Playback{
		def:     &discordgo.ApplicationCommand{},
		storage: &sync.Map{},
	}

	err := target.Handler(&discordgo.Session{}, inter)

	if err == nil {
		t.Fatal("should not be empty")
	}
	entry, ok := target.storage.Load("test")

	if !ok {
		t.Fatal("missing entry from map")
	}

	if entry.(*sync.Mutex) == nil {
		t.Fatal("missing mutex")
	}
}

func joinVoiceMock(sess *discordgo.Session, guildID, voiceID string, mute, deaf bool) (*discordgo.VoiceConnection, error) {
	return &discordgo.VoiceConnection{}, nil
}

func getGuildMock(sess *discordgo.Session, inter *discordgo.InteractionCreate) (*discordgo.Guild, error) {
	return &discordgo.Guild{
		ID: "test",
		VoiceStates: []*discordgo.VoiceState{
			{
				UserID: "user",
			},
		},
	}, nil
}

func sendRespMock(session *discordgo.Session, interaction *discordgo.InteractionCreate, msg string) {
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

package playback

import (
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

	got, err := New(testSyncMap, testDef)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

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

	if target.Data() == nil {
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

func TestHandlerWhenCalledCreatesMutex(t *testing.T) {
	t.Parallel()

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
		def:       &discordgo.ApplicationCommand{},
		storage:   &sync.Map{},
		voiceFunc: joinVoiceStub,
		guildFunc: getGuildStub,
		respFunc:  defferRespStub,
		editFunc:  editRespStub,
	}

	err := target.Handle(&discordgo.Session{}, inter)

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

func joinVoiceStub(_ *discordgo.Session, _, _ string, _, _ bool) (*discordgo.VoiceConnection, error) {
	return &discordgo.VoiceConnection{}, nil
}

func getGuildStub(_ *discordgo.Session, _ *discordgo.InteractionCreate) (*discordgo.Guild, error) {
	return &discordgo.Guild{
		ID: "test",
		VoiceStates: []*discordgo.VoiceState{
			{
				UserID: "user",
			},
		},
	}, nil
}

func defferRespStub(_ *discordgo.Session, _ *discordgo.InteractionCreate) error {
	return nil
}
func editRespStub(_ *discordgo.Session, _ *discordgo.InteractionCreate, _ string) error {
	return nil
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

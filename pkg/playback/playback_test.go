package playback

import (
	"os"
	"sync"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestDefinition(t *testing.T) {
	t.Parallel()

	target := &Playback{
		def: &discordgo.ApplicationCommand{},
	}

	if target.Definition() == nil {
		t.Fatal("This should not be nil!")
	}
}

func TestPlaySound(t *testing.T) {
	t.Parallel()

	testChan := make(chan []byte, 100)
	fh, _ := os.Open("testdata/test_file.dca")
	defer fh.Close()

	err := playSound(testChan, fh)

	if err != nil {
		t.Fatal("This should be nil!")
	}

	packet := <-testChan

	if packet == nil {
		t.Fatal("This should not be nil!")
	}
}

func TestPlaySound_WithError(t *testing.T) {
	t.Parallel()

	fh, _ := os.Open("")
	err := playSound(make(chan<- []byte), fh)

	if err == nil {
		t.Fatal("This should not be nil!")
	}
}

func TestGetAudioSource(t *testing.T) {
	tests := []struct {
		path       string
		shouldFail bool
	}{
		{"", true},
		{"testdata/test_file.dca", false},
		{"localhost.com/api", true},
		{"https://www.google.com/", false},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			res, err := getAudioSource(tc.path)

			if err != nil && !tc.shouldFail {
				t.Fatal(err)
			}

			if tc.shouldFail == false && res == nil {
				t.Fatal("Should not be nil!")
			}
		})
	}
}

func TestHandler(t *testing.T) {
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
		t.Fatal("Should not be emtty!")
	}
	entry, ok := target.storage.Load("test")

	if !ok {
		t.Fatal("Missing entry from map!")
	}

	if entry.(*sync.Mutex) == nil {
		t.Fatal("Missing mutex!")
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

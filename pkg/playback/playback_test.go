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

func TestPlaySoundsWithError(t *testing.T) {
	t.Parallel()

	fh, _ := os.Open("")
	err := playSound(make(chan<- []byte), fh)

	if err == nil {
		t.Error("this should not be nil")
	}
}

// TODO: Use a mock RoundTripFunc
func TestGetAudioSourceWithInvalidURI(t *testing.T) {
	t.Parallel()

	tests := []string{"", "localhost.com/api"}

	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			_, err := getAudioSource(tc)

			if err == nil {
				t.Error(err)
			}
		})
	}
}

func TestGetAudioSourceWithValidURI(t *testing.T) {
	t.Parallel()

	tests := []string{"testdata/test_file.dca", "https://www.google.com/"}

	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			res, err := getAudioSource(tc)

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

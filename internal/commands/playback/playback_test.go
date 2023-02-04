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
		t.Error("This should not be nil!")
	}
}

func TestPlaySoundWithFile(t *testing.T) {
	t.Parallel()

	testChan := make(chan []byte, 100)
	fh, _ := os.Open("testdata/test_file.dca")
	defer fh.Close()

	err := playSound(testChan, fh)

	if err != nil {
		t.Error("This should be nil!")
	}

	packet := <-testChan

	if packet == nil {
		t.Error("This should not be nil!")
	}
}

func TestPlaySound_WithError(t *testing.T) {
	t.Parallel()

	fh, _ := os.Open("")
	err := playSound(make(chan<- []byte), fh)

	if err == nil {
		t.Error("This should not be nil!")
	}
}

func TestIdleDisconnect(t *testing.T) {
	t.Parallel()
	target := cmdSync{}

	target.idleDisconnect(&discordgo.VoiceConnection{})

	if target.idle == nil {
		t.Error("This should not be nil!")
	}
}

func TestHandler(t *testing.T) {
	t.Parallel()
	voice = joinVoiceMock
	guild = getGuildMock
	response = sendRespMock

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
		t.Errorf("Should not be emtty!")
	}
	entry, ok := target.storage.Load("test")

	if !ok {
		t.Errorf("Missing entry from map!")
	}

	if entry.(*cmdSync).idle == nil {
		t.Errorf("Missing timer!")
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
	return
}

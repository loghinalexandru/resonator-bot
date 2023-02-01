package playback

import (
	"sync"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestDefinition(t *testing.T) {
	t.Parallel()

	target := &Playback{
		Def: &discordgo.ApplicationCommand{},
	}

	if target.Definition() == nil {
		t.Error("This should not be nil!")
	}
}

func TestPlaySoundWithFile(t *testing.T) {
	t.Parallel()

	testChan := make(chan []byte, 100)
	res := playSound(testChan, "mock/test_file.dca")

	if res != nil {
		t.Error("This should be nil!")
	}

	packet := <-testChan

	if packet == nil {
		t.Error("This should not be nil!")
	}
}

func TestPlaySoundWithError(t *testing.T) {
	t.Parallel()

	res := playSound(make(chan<- []byte), "")

	if res == nil {
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
	guildID := "test"
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
		Def:     &discordgo.ApplicationCommand{},
		Storage: &sync.Map{},
		voice: func(sess *discordgo.Session, guildID, voiceID string, mute, deaf bool) (*discordgo.VoiceConnection, error) {
			return &discordgo.VoiceConnection{}, nil
		},
		guild: func(sess *discordgo.Session, inter *discordgo.InteractionCreate) (*discordgo.Guild, error) {
			return &discordgo.Guild{
				ID: guildID,
				VoiceStates: []*discordgo.VoiceState{
					{
						UserID: "user",
					},
				},
			}, nil
		},
		response: func(sess *discordgo.Session, interaction *discordgo.InteractionCreate, msg string) {
			return
		},
	}

	err := target.Handler(&discordgo.Session{}, inter)

	if err == nil {
		t.Errorf("Should not be emtty!")
	}
	entry, ok := target.Storage.Load(guildID)

	if !ok {
		t.Errorf("Missing entry from map!")
	}

	if entry.(*cmdSync).idle == nil {
		t.Errorf("Missing timer!")
	}
}

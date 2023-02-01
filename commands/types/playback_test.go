package types

import (
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

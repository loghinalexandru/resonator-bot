package types

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

type Playback struct {
	Storage *sync.Map
	Def     *discordgo.ApplicationCommand
}

type cmdSync struct {
	mtx  sync.Mutex
	idle *time.Timer
}

func (cmd *Playback) Definition() *discordgo.ApplicationCommand {
	return cmd.Def
}

func (cmd *Playback) Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var botvc *discordgo.VoiceConnection
	var err error

	channel, _ := sess.State.Channel(inter.ChannelID)
	guild, _ := sess.State.Guild(channel.GuildID)

	for _, voice := range guild.VoiceStates {
		if inter.Member.User.ID == voice.UserID {
			botvc, err = sess.ChannelVoiceJoin(guild.ID, voice.ChannelID, false, true)
		}
	}

	if botvc == nil || err != nil {
		sendResponse(sess, inter, "Please join a voice channel!")
		return err
	}

	entry, _ := cmd.Storage.LoadOrStore(guild.ID, &cmdSync{})
	cmdSync := entry.(*cmdSync)
	result := cmdSync.mtx.TryLock()

	if !result {
		sendResponse(sess, inter, "Please wait your turn!")
		return nil
	}

	if cmdSync.idle != nil {
		cmdSync.idle.Stop()
		cmdSync.idle = nil
	}

	defer cmdSync.idleDisconnect(botvc)
	defer cmdSync.mtx.Unlock()
	defer botvc.Speaking(false)

	botvc.Speaking(true)
	sendResponse(sess, inter, "Playing!")

	path := fmt.Sprintf("%v", inter.ApplicationCommandData().Options[0].Value)
	err = playSound(botvc.OpusSend, path)

	if err != nil {
		return err
	}

	return nil
}

func playSound(soundBuff chan<- []byte, filePath string) error {
	input, ioError := os.Open(filePath)
	if ioError != nil {
		return ioError
	}

	defer input.Close()
	decoder := dca.NewDecoder(input)

	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		select {
		case soundBuff <- frame:
		case <-time.After(2 * time.Second):
			return errors.New("Timeout!")
		}
	}

	return nil
}

func sendResponse(session *discordgo.Session, interaction *discordgo.InteractionCreate, msg string) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (cmdSync *cmdSync) idleDisconnect(vc *discordgo.VoiceConnection) {
	cmdSync.idle = time.AfterFunc(3*time.Minute, func() { vc.Disconnect() })
}

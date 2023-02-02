package playback

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

const (
	waitTurn = "Please wait your turn!"
	joinVc   = "Please join a voice channel!"
	exec     = "Playing!"
)

var (
	voice    = joinVoice
	guild    = getGuild
	response = sendResp
)

type Playback struct {
	storage *sync.Map
	def     *discordgo.ApplicationCommand
}

type cmdSync struct {
	mtx  sync.Mutex
	idle *time.Timer
}

func (cmd *Playback) Definition() *discordgo.ApplicationCommand {
	return cmd.def
}

func (cmd *Playback) Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var botvc *discordgo.VoiceConnection
	var err error

	guild, _ := guild(sess, inter)

	for _, vc := range guild.VoiceStates {
		if inter.Member.User.ID == vc.UserID {
			botvc, err = voice(sess, guild.ID, vc.ChannelID, false, true)
		}
	}

	if botvc == nil || err != nil {
		response(sess, inter, joinVc)
		return err
	}

	entry, _ := cmd.storage.LoadOrStore(guild.ID, &cmdSync{})
	cmdSync := entry.(*cmdSync)
	result := cmdSync.mtx.TryLock()

	if !result {
		response(sess, inter, waitTurn)
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
	response(sess, inter, exec)

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

func joinVoice(sess *discordgo.Session, guildID, voiceID string, mute, deaf bool) (*discordgo.VoiceConnection, error) {
	return sess.ChannelVoiceJoin(guildID, voiceID, mute, deaf)
}

func getGuild(sess *discordgo.Session, inter *discordgo.InteractionCreate) (*discordgo.Guild, error) {
	channel, _ := sess.State.Channel(inter.ChannelID)
	return sess.State.Guild(channel.GuildID)
}

func sendResp(session *discordgo.Session, interaction *discordgo.InteractionCreate, msg string) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

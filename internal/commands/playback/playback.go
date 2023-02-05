package playback

import (
	"errors"
	"io"
	"net/http"
	"net/url"
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

	entry, _ := cmd.storage.LoadOrStore(guild.ID, &sync.Mutex{})
	mtx := entry.(*sync.Mutex)
	result := mtx.TryLock()

	if !result {
		response(sess, inter, waitTurn)
		return nil
	}

	defer mtx.Unlock()
	defer botvc.Speaking(false)

	botvc.Speaking(true)
	response(sess, inter, exec)

	var input io.Reader
	path := inter.ApplicationCommandData().Options[0].Value.(string)

	url, err := url.Parse(path)

	//TODO: refactor this & add error handling
	if url.Scheme == "" || url.Host == "" || url.Path == "" {
		fileReader, _ := os.Open(path)
		input = fileReader
		defer fileReader.Close()
	} else {
		resp, _ := http.Get(url.String())
		input = resp.Body
		defer resp.Body.Close()
	}

	err = playSound(botvc.OpusSend, input)

	if err != nil {
		return err
	}

	return nil
}

func playSound(soundBuff chan<- []byte, fh io.Reader) error {
	if fh == nil {
		return errors.New("Null file handler!")
	}

	decoder := dca.NewDecoder(fh)

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

// Seam functions for testing purposes
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

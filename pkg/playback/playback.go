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
	voice   = joinVoice
	guild   = getGuild
	respond = sendResp
)

type Playback struct {
	storage *sync.Map
	def     *discordgo.ApplicationCommand
}

func New(syncMap *sync.Map, definition *discordgo.ApplicationCommand) *Playback {
	return &Playback{
		def:     definition,
		storage: syncMap,
	}
}

func (cmd *Playback) Definition() *discordgo.ApplicationCommand {
	return cmd.def
}

func (cmd *Playback) Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var err error

	guild, _ := guild(sess, inter)
	botvc, ok := sess.VoiceConnections[guild.ID]

	if !ok {
		for _, vc := range guild.VoiceStates {
			if inter.Member.User.ID == vc.UserID {
				botvc, err = voice(sess, guild.ID, vc.ChannelID, false, true)
			}
		}
	}

	if botvc == nil || err != nil {
		respond(sess, inter, joinVc)
		return err
	}

	entry, _ := cmd.storage.LoadOrStore(guild.ID, &sync.Mutex{})
	mtx := entry.(*sync.Mutex)
	result := mtx.TryLock()

	if !result {
		respond(sess, inter, waitTurn)
		return nil
	}

	defer mtx.Unlock()
	defer botvc.Speaking(false)

	botvc.Speaking(true)
	respond(sess, inter, exec)

	path := inter.ApplicationCommandData().Options[0].Value.(string)
	input, err := getAudioSource(path)

	if err != nil {
		return err
	}

	defer input.Close()
	err = playSound(botvc.OpusSend, input)

	if err != nil {
		return err
	}

	return nil
}

func getAudioSource(path string) (io.ReadCloser, error) {
	url, err := url.Parse(path)

	if err != nil || url.Scheme == "" || url.Host == "" {
		res, err := os.Open(path)

		if err != nil {
			return nil, err
		}

		return res, nil
	}

	res, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func playSound(soundBuff chan<- []byte, fh io.Reader) error {
	if fh == nil {
		return errors.New("null file handler")
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
			return errors.New("timeout")
		}
	}

	return nil
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

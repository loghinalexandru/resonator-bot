package playback

import (
	"errors"
	"io"
	"net/http"
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

type playbackOpt func(*Playback)

type Playback struct {
	storage *sync.Map
	url     string
	def     *discordgo.ApplicationCommand
}

func New(syncMap *sync.Map, definition *discordgo.ApplicationCommand, opts ...playbackOpt) *Playback {
	result := &Playback{
		def:     definition,
		storage: syncMap,
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

func WithURL(remoteURL string) playbackOpt {
	return func(p *Playback) {
		p.url = remoteURL
	}
}

func (cmd *Playback) Definition() *discordgo.ApplicationCommand {
	return cmd.def
}

func (cmd *Playback) Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var err error

	guild, _ := guild(sess, inter)
	botvc, exists := sess.VoiceConnections[guild.ID]

	if !exists {
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
	success := mtx.TryLock()

	if !success {
		respond(sess, inter, waitTurn)
		return nil
	}

	defer mtx.Unlock()
	defer botvc.Speaking(false)

	botvc.Speaking(true)

	path := inter.ApplicationCommandData().Options[0].Value.(string)
	input, err := getAudioSource(cmd.url, path, http.DefaultClient)

	if err != nil {
		respond(sess, inter, err.Error())
		return err
	}

	//Need to respond faster somehow
	respond(sess, inter, exec)
	err = playSound(botvc.OpusSend, input)

	if err != nil {
		return err
	}

	return nil
}

func getAudioSource(baseURL string, path string, client *http.Client) (io.ReadCloser, error) {
	if baseURL != "" {
		res, err := client.Get(baseURL + path)

		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			return nil, errors.New("could not retreive data from provided URL")
		}

		return res.Body, nil
	}

	res, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func playSound(soundBuff chan<- []byte, fh io.ReadCloser) error {
	if fh == nil {
		return errors.New("nil file handler")
	}

	defer fh.Close()
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
			return errors.New("opus frame timeout")
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

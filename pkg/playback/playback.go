package playback

import (
	"errors"
	"fmt"
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
	msgMissingVoiceChannel = "Please join a voice channel first!"
	msgConcurrentPlayback  = "Please wait your turn!"
	msgMissingAudio        = "Could not retrieve specified audio!"
	msgSuccess             = "Playback started!"
)

var (
	ErrFileHandler = errors.New("nil file handler")
	ErrTimeout     = errors.New("opus frame timeout")
	ErrContent     = errors.New("could not retreive data from provided URL")
	voice          = joinVoice
	guild          = getGuild
	respond        = sendResp
)

type playbackOpt func(*Playback) error

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
		//Handle error case
		opt(result)
	}

	return result
}

func WithURL(remoteURL string) playbackOpt {
	return func(p *Playback) error {
		_, err := url.Parse(remoteURL)

		if err != nil {
			return err
		}

		p.url = remoteURL
		return nil
	}
}

func (cmd *Playback) Definition() *discordgo.ApplicationCommand {
	return cmd.def
}

func (cmd *Playback) Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) (err error) {
	guild, _ := guild(sess, inter)
	botvc, exists := sess.VoiceConnections[guild.ID]

	if !exists || inter.Member.User.ID != botvc.UserID {
		for _, vc := range guild.VoiceStates {
			if inter.Member.User.ID == vc.UserID {
				botvc, err = voice(sess, guild.ID, vc.ChannelID, false, true)
			}
		}
	}

	if botvc == nil || err != nil {
		respond(sess, inter, msgMissingVoiceChannel)
		return err
	}

	entry, _ := cmd.storage.LoadOrStore(guild.ID, &sync.Mutex{})
	mtx := entry.(*sync.Mutex)
	ok := mtx.TryLock()

	if !ok {
		respond(sess, inter, msgConcurrentPlayback)
		return nil
	}

	botvc.Speaking(true)

	defer mtx.Unlock()
	defer func() {
		err = errors.Join(err, botvc.Speaking(false))
	}()

	userOpt := inter.ApplicationCommandData().Options[0].Value.(string)
	var audio io.ReadCloser

	if cmd.url != "" {
		audio, err = remoteAudioSource(fmt.Sprintf(cmd.url, userOpt), http.DefaultClient)
	} else {
		audio, err = localAudioSource(userOpt)
	}

	if err != nil {
		respond(sess, inter, msgMissingAudio)
		return err
	}

	respond(sess, inter, msgSuccess)
	err = playSound(botvc.OpusSend, audio)

	if err != nil {
		return err
	}

	return nil
}

func remoteAudioSource(remoteURL string, client *http.Client) (io.ReadCloser, error) {
	res, err := client.Get(remoteURL)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrContent
	}

	return res.Body, nil
}

func localAudioSource(path string) (io.ReadCloser, error) {
	res, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func playSound(soundBuff chan<- []byte, fh io.ReadCloser) error {
	if fh == nil {
		return ErrFileHandler
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
			return ErrTimeout
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

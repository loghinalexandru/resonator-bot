package playback

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/loghinalexandru/resonator/pkg/audio"
)

const (
	msgMissingVoiceChannel = "Please join a voice channel first!"
	msgErrorOnSpeak        = "Could not unmute bot!"
	msgConcurrentPlayback  = "Please wait your turn!"
	msgMissingAudio        = "Could not retrieve specified audio!"
	msgSuccess             = "Playback started!"
)

var (
	ErrFileHandler = errors.New("nil file handler")
	ErrTimeout     = errors.New("opus frame timeout")
	ErrHttpClient  = errors.New("missing http client")
	ErrAudioSource = errors.New("missing audio source")
	voice          = joinVoice
	guild          = getGuild
	respond        = sendResp
)

type playbackOpt func(*Playback) error

type AudioProvider interface {
	Audio(path string) (io.ReadCloser, error)
}

type Playback struct {
	source  AudioProvider
	storage *sync.Map
	def     *discordgo.ApplicationCommand
}

func New(syncMap *sync.Map, definition *discordgo.ApplicationCommand, opts ...playbackOpt) *Playback {
	result := &Playback{
		def:     definition,
		source:  audio.Local{},
		storage: syncMap,
	}

	for _, opt := range opts {
		//Handle error case
		opt(result)
	}

	return result
}

func WithAudioSource(source AudioProvider) playbackOpt {
	return func(p *Playback) error {
		if source == nil {
			return ErrAudioSource
		}

		p.source = source
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

	err = botvc.Speaking(true)

	if err != nil {
		respond(sess, inter, msgErrorOnSpeak)
		return err
	}

	defer mtx.Unlock()
	defer func() {
		err = errors.Join(err, botvc.Speaking(false))
	}()

	userOpt := inter.ApplicationCommandData().Options[0].Value.(string)
	audio, err := cmd.source.Audio(userOpt)

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

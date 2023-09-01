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
	voice          = discordgoVoice
	guild          = discordgoGuild
	defferResponse = discordgoInteractionResp
	respond        = discordgoInteractionEdit
)

type playbackOpt func(*Playback) error

type Playback struct {
	src     audio.Provider
	storage *sync.Map
	def     *discordgo.ApplicationCommand
}

func New(syncMap *sync.Map, definition *discordgo.ApplicationCommand, opts ...playbackOpt) (*Playback, error) {
	result := &Playback{
		def:     definition,
		src:     audio.NewLocal(),
		storage: syncMap,
	}

	for _, opt := range opts {
		err := opt(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func WithSource(provider audio.Provider) playbackOpt {
	return func(p *Playback) error {
		if provider == nil {
			return ErrAudioSource
		}

		p.src = provider
		return nil
	}
}

func (cmd *Playback) Data() *discordgo.ApplicationCommand {
	return cmd.def
}

func (command *Playback) Handle(sess *discordgo.Session, inter *discordgo.InteractionCreate) (err error) {
	err = defferResponse(sess, inter)

	if err != nil {
		return err
	}

	guild, _ := guild(sess, inter)
	entry, _ := command.storage.LoadOrStore(guild.ID, &sync.Mutex{})
	mtx := entry.(*sync.Mutex)
	ok := mtx.TryLock()

	if !ok {
		err = respond(sess, inter, msgConcurrentPlayback)
		if err != nil {
			return err
		}

		return nil
	}

	defer mtx.Unlock()

	botvc, exists := sess.VoiceConnections[guild.ID]

	if !exists || inter.Member.User.ID != botvc.UserID {
		for _, vc := range guild.VoiceStates {
			if inter.Member.User.ID == vc.UserID {
				botvc, err = voice(sess, guild.ID, vc.ChannelID, false, true)
			}
		}
	}

	if botvc == nil || err != nil {
		return errors.Join(err, respond(sess, inter, msgMissingVoiceChannel))
	}

	err = botvc.Speaking(true)

	if err != nil {
		return errors.Join(err, respond(sess, inter, msgErrorOnSpeak))
	}

	defer func() {
		err = errors.Join(err, botvc.Speaking(false))
	}()

	userOpt := inter.ApplicationCommandData().Options[0].Value.(string)
	audio, err := command.src.Audio(userOpt)

	if err != nil {
		return errors.Join(err, respond(sess, inter, msgMissingAudio))
	}

	err = respond(sess, inter, msgSuccess)
	if err != nil {
		return err
	}

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
func discordgoVoice(sess *discordgo.Session, guildID, voiceID string, mute, deaf bool) (*discordgo.VoiceConnection, error) {
	return sess.ChannelVoiceJoin(guildID, voiceID, mute, deaf)
}

func discordgoGuild(sess *discordgo.Session, inter *discordgo.InteractionCreate) (*discordgo.Guild, error) {
	channel, _ := sess.State.Channel(inter.ChannelID)
	return sess.State.Guild(channel.GuildID)
}

func discordgoInteractionResp(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

func discordgoInteractionEdit(session *discordgo.Session, interaction *discordgo.InteractionCreate, msg string) error {
	_, err := session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})

	return err
}

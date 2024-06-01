package playback

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/loghinalexandru/resonator/pkg/provider"
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
	ErrAudioSource = errors.New("missing audio source")
)

type Opt func(*Playback) error

type voice func(sess *discordgo.Session, guildID, voiceID string, mute, deaf bool) (*discordgo.VoiceConnection, error)
type guild func(sess *discordgo.Session, inter *discordgo.InteractionCreate) (*discordgo.Guild, error)
type interactionResp func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error
type interactionEdit func(session *discordgo.Session, interaction *discordgo.InteractionCreate, msg string) error

type Source interface {
	Fetch(path string) (io.ReadCloser, error)
}

type Playback struct {
	src       Source
	storage   *sync.Map
	def       *discordgo.ApplicationCommand
	voiceFunc voice
	guildFunc guild
	respFunc  interactionResp
	editFunc  interactionEdit
}

func New(syncMap *sync.Map, definition *discordgo.ApplicationCommand, opts ...Opt) (*Playback, error) {
	return newInternal(syncMap, definition, discordgoVoice, discordgoGuild, discordgoInteractionResp, discordgoInteractionEdit, opts...)
}

func WithSource(provider Source) Opt {
	return func(p *Playback) error {
		if provider == nil {
			return ErrAudioSource
		}

		p.src = provider
		return nil
	}
}

func newInternal(syncMap *sync.Map, definition *discordgo.ApplicationCommand, voice voice, guild guild, resp interactionResp, edit interactionEdit, opts ...Opt) (*Playback, error) {
	result := &Playback{
		def:       definition,
		src:       &provider.LocalProvider{},
		storage:   syncMap,
		voiceFunc: voice,
		guildFunc: guild,
		respFunc:  resp,
		editFunc:  edit,
	}

	for _, opt := range opts {
		err := opt(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (command *Playback) Data() *discordgo.ApplicationCommand {
	return command.def
}

func (command *Playback) Handle(sess *discordgo.Session, inter *discordgo.InteractionCreate) (err error) {
	err = command.respFunc(sess, inter)

	if err != nil {
		return err
	}

	guild, _ := command.guildFunc(sess, inter)
	entry, _ := command.storage.LoadOrStore(guild.ID, &sync.Mutex{})
	mtx := entry.(*sync.Mutex)
	ok := mtx.TryLock()

	if !ok {
		err = command.editFunc(sess, inter, msgConcurrentPlayback)
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
				botvc, err = command.voiceFunc(sess, guild.ID, vc.ChannelID, false, true)
			}
		}
	}

	if botvc == nil || err != nil {
		return errors.Join(err, command.editFunc(sess, inter, msgMissingVoiceChannel))
	}

	err = botvc.Speaking(true)

	if err != nil {
		return errors.Join(err, command.editFunc(sess, inter, msgErrorOnSpeak))
	}

	defer func() {
		err = errors.Join(err, botvc.Speaking(false))
	}()

	userOpt := inter.ApplicationCommandData().Options[0].Value.(string)
	rawAudio, err := command.src.Fetch(userOpt)

	if err != nil {
		return errors.Join(err, command.editFunc(sess, inter, msgMissingAudio))
	}

	err = command.editFunc(sess, inter, msgSuccess)
	if err != nil {
		return err
	}

	err = playSound(botvc.OpusSend, rawAudio)
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

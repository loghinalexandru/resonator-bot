package commands

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type playCommand struct {
	identifier string
}

func (cmd playCommand) ID() string {
	return cmd.identifier
}

func (cmd playCommand) Definition() *discordgo.ApplicationCommand {
	result := new(discordgo.ApplicationCommand)
	result.Name = cmd.identifier
	result.Description = "This command is used to play a sound in the chat!"
	result.Options = append(result.Options, &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "type",
		Description: "Sound type to be played!",
		Required:    true,
	})

	return result
}

func (playCommand) Handler(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	channel, _ := session.State.Channel(interaction.ChannelID)
	guild, _ := session.State.Guild(channel.GuildID)

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Playing!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	for _, voice := range guild.VoiceStates {
		if interaction.Member.User.ID == voice.UserID {
			botvc, error := session.ChannelVoiceJoin(guild.ID, voice.ChannelID, false, true)

			if error != nil {
				return error
			}

			botvc.Speaking(true)

			path := fmt.Sprintln("misc/", interaction.ApplicationCommandData().Options[0].Value, ".dca")
			soundError := play(botvc, path)

			if soundError != nil {
				return soundError
			}

			botvc.Speaking(false)
		}
	}

	return nil
}

func play(voice *discordgo.VoiceConnection, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	var opuslen int16
	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			return err
		}

		opusFrame := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &opusFrame)

		if err != nil {
			return err
		}

		select {
		case voice.OpusSend <- opusFrame:
		case <-time.After(2 * time.Second):
			return errors.New("Timeout!")
		}
	}
}

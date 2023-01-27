package commands

import (
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

type playCommand struct {
	identifier string
	mutex      *sync.Mutex
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

func (cmd playCommand) Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error {
	channel, _ := sess.State.Channel(inter.ChannelID)
	guild, _ := sess.State.Guild(channel.GuildID)

	result := cmd.mutex.TryLock()

	if !result {
		sendResponse(sess, inter, "Please wait your turn!")
		return nil
	}
	sendResponse(sess, inter, "Playing!")

	defer cmd.mutex.Unlock()

	for _, voice := range guild.VoiceStates {
		if inter.Member.User.ID == voice.UserID {
			botvc, err := sess.ChannelVoiceJoin(guild.ID, voice.ChannelID, false, true)

			if err != nil {
				return err
			}

			path := "misc/" + inter.ApplicationCommandData().Options[0].Value.(string) + ".dca"
			err = playSound(sess, botvc, path)

			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func playSound(sess *discordgo.Session, voice *discordgo.VoiceConnection, filePath string) error {
	voice.Speaking(true)
	defer voice.Speaking(false)

	input, ioError := os.Open(filePath)
	if ioError != nil {
		return ioError
	}

	defer input.Close()
	decoder := dca.NewDecoder(input)

	if !voice.Ready {
		return errors.New("Voice channel not ready!")
	}

	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		select {
		case voice.OpusSend <- frame:
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

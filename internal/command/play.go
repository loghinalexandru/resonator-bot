package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/bot"
	"github.com/loghinalexandru/resonator/pkg/playback"
)

func newPlay(ctx *bot.Context) *playback.Playback {
	result, err := playback.New(ctx.Sync, &discordgo.ApplicationCommand{
		Name:        "play",
		Description: "This command is used to play a sound in the chat!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "type",
				Description: "Sound type to be played!",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Ara-Ara",
						Value: "misc/ara.dca",
					},
					{
						Name:  "Yoooooooouuu",
						Value: "misc/yoo.dca",
					},
					{
						Name:  "FBI Open Up",
						Value: "misc/fbi.dca",
					},
					{
						Name:  "Hehe Boy!",
						Value: "misc/heheboy.dca",
					},
					{
						Name:  "Yamete Kudasai!",
						Value: "misc/yamete.dca",
					},
					{
						Name:  "No God Please No!",
						Value: "misc/nogod.dca",
					},
					{
						Name:  "Mission failed.",
						Value: "misc/failed.dca",
					},
					{
						Name:  "Death",
						Value: "misc/death.dca",
					},
					{
						Name:  "UwU",
						Value: "misc/uwu.dca",
					},
					{
						Name:  "Fast AF",
						Value: "misc/fastaf.dca",
					},
					{
						Name:  "City Boy",
						Value: "misc/cityboy.dca",
					},
					{
						Name:  "So good",
						Value: "misc/sogood.dca",
					},
					{
						Name:  "Why are u running",
						Value: "misc/whyurunning.dca",
					},
					{
						Name:  "Aw shit",
						Value: "misc/awshit.dca",
					},
					{
						Name:  "Fuck fuck",
						Value: "misc/fuck.dca",
					},
					{
						Name:  "Sitcom laugh",
						Value: "misc/sitcomlaugh.dca",
					},
					{
						Name:  "Law & Order",
						Value: "misc/lawandorder.dca",
					},
					{
						Name:  "What the fuck",
						Value: "misc/whatthefuck.dca",
					},
					{
						Name:  "Bloody fuck you",
						Value: "misc/bloodyfuckyou.dca",
					},
					{
						Name:  "I don't want peace!",
						Value: "misc/iwantproblems.dca",
					},
					{
						Name:  "Why are these fools",
						Value: "misc/thesefools.dca",
					},
				},
				Required: true,
			},
		},
	})

	if err != nil {
		ctx.Logger.Error("Error creating /play command", "err", err)
	}

	return result
}

package playback

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

func NewCurse(sync *sync.Map) *Playback {
	out := Playback{
		storage: sync,
		def: &discordgo.ApplicationCommand{
			Name:        "curse",
			Description: "This command is used to play a friendly encouragement!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "language",
					Description: "Specify in which language your encouragement will be!",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Romanian",
							//Reference this via ENV variable
							Value: "http://swears-svc/api/random/file?lang=ro&opus=true",
						},
					},
					Required: true,
				},
			},
		},
	}

	return &out
}

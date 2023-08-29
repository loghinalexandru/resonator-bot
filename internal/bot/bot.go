package bot

import (
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Context struct {
	Sync      *sync.Map
	SwearsAPI *url.URL
	Logger    Logger
}

func NewContext() *Context {
	ctx := &Context{
		Sync:      &sync.Map{},
		SwearsAPI: swearsAPI("SWEARS_API_URL"),
		Logger:    logWithLvl("LOG_LEVEL"),
	}

	return ctx
}

func swearsAPI(varName string) *url.URL {
	res, err := url.Parse(os.Getenv(varName))
	if err != nil {
		panic(err)
	}

	return res
}

func Token() string {
	return os.Getenv("BOT_TOKEN")
}

func Intents() discordgo.Intent {
	return discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
}

func ID() int {
	id := os.Getenv("SHARD_ID")
	if id == "" {
		return 0
	}

	index := strings.Split(id, "-")
	res, err := strconv.Atoi(index[len(index)-1])

	if err != nil {
		panic(err)
	}

	return res
}

func Shards() int {
	replicas := os.Getenv("SHARD_COUNT")
	if replicas == "" {
		return 1
	}

	res, err := strconv.Atoi(replicas)

	if err != nil {
		panic(err)
	}

	return res
}

func Register(sess *discordgo.Session, ctx *Context) {
	sess.AddHandler(func(sess *discordgo.Session, gld *discordgo.GuildCreate) {
		ctx.Logger.Info("joined guild", "guildID", gld.ID)
	})
}

package bot

import (
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	Req       prometheus.Counter
	Err       prometheus.Counter
}

func NewContext() *Context {
	ctx := &Context{
		Sync:      &sync.Map{},
		SwearsAPI: envToURL("SWEARS_API_URL"),
		Logger:    logWithLvl("LOG_LEVEL"),
		Req: promauto.NewCounter(prometheus.CounterOpts{
			Name: "resonator_commands_invoked_total",
			Help: "The total number of commands invoked",
		}),
		Err: promauto.NewCounter(prometheus.CounterOpts{
			Name: "resonator_command_errors_total",
			Help: "The total number of commands errors",
		}),
	}

	return ctx
}

func envToURL(varName string) *url.URL {
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

func Cleanup() bool {
	c := os.Getenv("DEREGISTER_COMMANDS")
	if c == "" {
		return false
	}

	res, err := strconv.ParseBool(c)

	if err != nil {
		panic(err)
	}

	return res
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

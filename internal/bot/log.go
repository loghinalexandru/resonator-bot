package bot

import (
	"log/slog"
	"os"
	"strings"
)

func logWithLvl(varName string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     strToLogLvl(os.Getenv(varName)),
	}))
}

func strToLogLvl(lvl string) slog.Level {
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}

}

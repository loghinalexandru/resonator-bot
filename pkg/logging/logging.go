package logging

import (
	"fmt"
	"log"
	"strings"
)

type LogLevel int

const (
	Debug LogLevel = iota + 1
	Info
	Warning
	Error
)

type Logger struct {
	logger *log.Logger
	minLvl LogLevel
}

func New(lvl LogLevel, logger *log.Logger) *Logger {
	return &Logger{
		minLvl: lvl,
		logger: logger,
	}
}

func (log *Logger) logInternal(lvl LogLevel, args ...any) {
	if lvl >= log.minLvl {
		log.logger.Printf(fmt.Sprintf("%v: %v", ToStr(lvl), fmt.Sprint(args...)))
	}
}

func (log *Logger) Debug(args ...any) {
	log.logInternal(Debug, args...)
}

func (log *Logger) Info(args ...any) {
	log.logInternal(Info, args...)
}

func (log *Logger) Warning(args ...any) {
	log.logInternal(Warning, args...)
}

func (log *Logger) Error(args ...any) {
	log.logInternal(Error, args...)
}

func ToStr(lvl LogLevel) string {
	switch lvl {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	default:
		panic("invalid value")
	}
}

func ToLogLevel(lvl string) LogLevel {
	sanitizedLvl := strings.ToUpper(strings.TrimSpace(lvl))

	switch sanitizedLvl {
	case "DEBUG":
		return Debug
	case "INFO":
		return Info
	case "WARNING":
		return Warning
	case "ERROR":
		return Error
	default:
		panic("invalid value")
	}
}

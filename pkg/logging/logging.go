package logging

import (
	"fmt"
	"log"
)

type logLevel int

const (
	Debug logLevel = iota + 1
	Info
	Warning
	Error
)

type Logger struct {
	logger *log.Logger
	minLvl logLevel
}

func New(lvl logLevel, logger *log.Logger) *Logger {
	return &Logger{
		minLvl: lvl,
		logger: logger,
	}
}

func (log *Logger) logInternal(lvl logLevel, args ...any) {
	if lvl >= log.minLvl {
		log.logger.Println(fmt.Sprintf("%v: %v", strLogLevel(lvl), fmt.Sprint(args...)))
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

func strLogLevel(lvl logLevel) string {
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
		panic("Unknown log level!")
	}
}

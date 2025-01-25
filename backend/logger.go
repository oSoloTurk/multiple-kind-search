package main

import (
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

// InitLogger initializes the global logger
func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	logger = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Caller().
		Logger()
}

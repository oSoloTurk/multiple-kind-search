package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger is the global logger instance
var Logger zerolog.Logger

// InitLogger initializes the global logger
func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}
	Logger = zerolog.New(consoleWriter).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Caller().
		Logger()
}

// SetLogLevel sets the global log level
func SetLogLevel(level zerolog.Level) {
	Logger = Logger.Level(level)
}

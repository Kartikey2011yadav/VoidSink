package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(level string, logFile string, logFormat string) {
	// Set global log level
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)

	var output io.Writer = os.Stderr

	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Error().Err(err).Msg("Failed to open log file, using stderr")
		} else {
			output = f
		}
	}

	if logFormat == "json" {
		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	} else {
		// Pretty print for console
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: output, TimeFormat: time.RFC3339})
	}
}

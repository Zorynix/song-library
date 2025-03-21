package app

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func SetupLogger(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		Logger.Error().Err(err).Msg("Invalid log level, defaulting to Info")
		return
	}

	zerolog.SetGlobalLevel(lvl)
	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	Logger.Info().Msgf("Logger initialized with level: %s", lvl.String())
}

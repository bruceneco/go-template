package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger(cfg *EnvConfig) zerolog.Logger {
	switch cfg.GoEnv {
	case EnvTypeDevelopment:
		log.Logger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	case EnvTypeProduction, EnvTypeStaging:
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	return log.Logger
}

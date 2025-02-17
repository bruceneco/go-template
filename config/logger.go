package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger(cfg *EnvConfig) {
	if cfg == nil {
		log.Fatal().Msg("config is nil")
		return
	}

	switch cfg.GoEnv {
	case EnvTypeProduction, EnvTypeStaging:
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger().Level(cfg.LogLevel)
	case EnvTypeDevelopment, EnvTypeTest:
		log.Logger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger().Level(cfg.LogLevel)
	}
}

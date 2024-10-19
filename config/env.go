package config

import (
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type EnvConfig struct {
	GoEnv         EnvType `env:"GO_ENV" envDefault:"development"`
	PostgresDSN   string  `env:"POSTGRES_DSN,required"`
	DBAutoMigrate bool    `env:"DB_AUTO_MIGRATE" envDefault:"false"`
	AMQPHost      string  `env:"AMQP_HOST" envDefault:"amqp://guest:guest@localhost"`
	HTTPPort      string  `env:"HTTP_PORT" envDefault:"3000"`
}

func LoadEnv() *EnvConfig {
	err := godotenv.Load(".env.local")
	if err != nil && os.IsNotExist(err) && isDevMode() {
		log.Fatal().Err(err).Msg("Failed to setup local environment variables")
	}

	var cfg EnvConfig
	if err = env.Parse(&cfg); err != nil {
		log.Panic().Err(err).Msg("failed to parse env variables")
	}
	return &cfg
}

func isDevMode() bool {
	envMode := os.Getenv("GO_ENV")
	return envMode == "development" || envMode == ""
}

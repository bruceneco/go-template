package config

import (
	"github.com/rs/zerolog"
	"os"
	"path"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type EnvConfig struct {
	GoEnv         EnvType       `env:"GO_ENV" envDefault:"development"`
	PostgresDSN   string        `env:"POSTGRES_DSN,required"`
	DBAutoMigrate bool          `env:"DB_AUTO_MIGRATE" envDefault:"false"`
	AMQPHost      string        `env:"AMQP_HOST" envDefault:"amqp://guest:guest@localhost"`
	HTTPPort      string        `env:"HTTP_PORT" envDefault:"3000"`
	LogLevel      zerolog.Level `env:"LOG_LEVEL" envDefault:"0"`
	ProjectRoot   string
}

func LoadEnv() *EnvConfig {
	filename := ".env"
	switch os.Getenv("GO_ENV") {
	case EnvTypeDevelopment.String():
		filename += ".local"
	case EnvTypeTest.String():
		filename += ".test"
	}
	root, err := findProjectRoot()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to find project root")
	}
	err = godotenv.Load(path.Join(root, filename))
	if err != nil && os.IsNotExist(err) && isDevMode() {
		log.Fatal().Err(err).Msg("Failed to setup local environment variables")
	}

	var cfg EnvConfig
	if err = env.Parse(&cfg); err != nil {
		log.Panic().Err(err).Msg("failed to parse env variables")
	}
	cfg.ProjectRoot = root
	return &cfg
}

func isDevMode() bool {
	envMode := os.Getenv("GO_ENV")
	return envMode == "development" || envMode == ""
}
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err = os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return "", os.ErrNotExist
}

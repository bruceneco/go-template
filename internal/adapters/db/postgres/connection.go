package postgres

import (
	"github.com/bruceneco/go-template/config"
	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm/logger"
	"path"
	"time"

	"github.com/rs/zerolog/log"
	psql "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	*gorm.DB
}

func NewConnection(cfg *config.EnvConfig) *Connection {
	db, err := gorm.Open(psql.Open(cfg.PostgresDSN), &gorm.Config{
		Logger: newGormLogger(),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres db")
	}
	if cfg.DBAutoMigrate {
		runMigrations(cfg, db)
	}
	return &Connection{db}
}

func runMigrations(cfg *config.EnvConfig, db *gorm.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: path.Join(cfg.ProjectRoot + "/tools/db/migrations"),
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get sql db")
	}
	n, err := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}
	log.Info().Msgf("Applied %d migrations", n)
}

func newGormLogger() logger.Interface {
	return logger.New(
		&log.Logger,
		logger.Config{
			SlowThreshold:             time.Second,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)
}

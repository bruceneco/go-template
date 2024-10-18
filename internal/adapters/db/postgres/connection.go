package postgres

import (
	"go-template/config"
	"go-template/internal/adapters/db/entities"

	"github.com/rs/zerolog/log"
	psql "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.EnvConfig) *gorm.DB {
	db, err := gorm.Open(psql.Open(cfg.PostgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres db")
	}
	migrate(db, cfg)
	return db
}

func migrate(db *gorm.DB, cfg *config.EnvConfig) {
	log.Info().Interface("cfg", cfg)
	if cfg.DBAutoMigrate {
		err := db.AutoMigrate(
			new(entities.Example),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to migrate db schema")
		}
	}
}

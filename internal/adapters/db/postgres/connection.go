package postgres

import (
	"go-template/config"
	"go-template/internal/adapters/db/entities"

	"github.com/rs/zerolog/log"
	psql "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	DB *gorm.DB
}

func NewConnection(cfg *config.EnvConfig) *Connection {
	db, err := gorm.Open(psql.Open(cfg.PostgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres db")
	}
	migrate(db, cfg)
	return &Connection{DB: db}
}

func migrate(db *gorm.DB, cfg *config.EnvConfig) {
	if cfg.DBAutoMigrate {
		err := db.AutoMigrate(
			new(entities.Example),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to migrate db schema")
		}
	}
}

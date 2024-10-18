package main

import (
	"go-template/config"
	"go-template/internal/adapters/db"
	"go-template/internal/adapters/db/entities"

	"github.com/ipfans/fxlogger"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"go.uber.org/fx"
)

func main() {
	env := config.LoadEnv()
	config.SetupLogger(env)
	fx.New(
		fx.Provide(func() *config.EnvConfig { return env }),
		fx.WithLogger(fxlogger.WithZerolog(log.Logger)),
		db.Module,
		fx.Invoke(func(db *gorm.DB) {
			err := db.Transaction(func(tx *gorm.DB) error {
				for range 3 {
					if err := tx.Create(&entities.Example{}).Error; err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create examples in database")
			}
			var examples []entities.Example
			if err = db.Find(&examples).Error; err != nil {
				log.Fatal().Err(err).Msg("failed to find examples in database")
			}
			log.Info().Int("num_entrie", len(examples)).Send()
		}),
	).Run()
}

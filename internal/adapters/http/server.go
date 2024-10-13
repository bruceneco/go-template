package http

import (
	"context"
	"go-template/config"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func NewHTTPServer(cfg *config.EnvConfig, lc fx.Lifecycle) *fiber.App {
	maxReadDuration := 5
	app := fiber.New(fiber.Config{
		Immutable:   true,
		ReadTimeout: time.Duration(maxReadDuration) * time.Second,
	})
	app.Use(fiberzerolog.New(fiberzerolog.Config{Logger: &log.Logger}))
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				err := app.Listen(":" + cfg.HTTPPort)
				if err != nil {
					log.Fatal().Err(err).Send()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		},
	})
	return app
}

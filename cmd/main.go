package main

import (
	"context"
	"go-template/config"
	"time"

	"github.com/ipfans/fxlogger"
	"github.com/rs/zerolog/log"

	"go.uber.org/fx"
)

func main() {
	env := config.LoadEnv()
	config.SetupLogger(env)
	fx.New(
		fx.Provide(func() *config.EnvConfig { return env }),
		fx.WithLogger(fxlogger.WithZerolog(log.Logger)),
		fx.Invoke(func(envConfig *config.EnvConfig) {
			log.Info().Str("env_mode", envConfig.GoEnv.String()).Msgf("env loaded")
		}),
		fx.Invoke(NewExample),
	).Run()
}

func NewExample(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info().Msg("Starting example")
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info().Msg("Stopping example")
			gracefulDelay := 3
			time.Sleep(time.Duration(gracefulDelay) * time.Second)
			log.Info().Msg("Stopped example")
			return nil
		},
	})
}

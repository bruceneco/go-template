package main

import (
	"github.com/ipfans/fxlogger"
	"github.com/rs/zerolog/log"
	"go-template/config"
	"go-template/internal/adapters"
	"go.uber.org/fx"
)

func main() {
	env := config.LoadEnv()
	config.SetupLogger(env)
	fx.New(
		fx.Provide(func() *config.EnvConfig { return env }),
		fx.WithLogger(fxlogger.WithZerolog(log.Logger)),
		adapters.Module,
	).Run()
}

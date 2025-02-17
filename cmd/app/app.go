package app

import (
	"github.com/ipfans/fxlogger"
	"github.com/rs/zerolog/log"
	"go-template/config"
	"go-template/internal/adapters"
	"go-template/internal/domain"
	"go.uber.org/fx"
)

func Inject() fx.Option {
	env := config.LoadEnv()
	config.SetupLogger(env)

	return fx.Options(
		fx.Provide(func() *config.EnvConfig { return env }),
		fx.WithLogger(fxlogger.WithZerolog(log.Logger)),
		adapters.Module,
		domain.Module,
	)
}

func Start() {
	fx.New(
		Inject(),
	).Run()
}

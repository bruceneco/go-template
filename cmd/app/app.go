package app

import (
	"github.com/bruceneco/go-template/config"
	"github.com/bruceneco/go-template/internal/adapters"
	"github.com/bruceneco/go-template/internal/adapters/grpc"
	"github.com/bruceneco/go-template/internal/adapters/http"
	"github.com/bruceneco/go-template/internal/domain"
	"github.com/ipfans/fxlogger"
	"github.com/rs/zerolog/log"
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
		fx.Invoke(http.Serve),
		fx.Invoke(grpc.Serve),
	).Run()
}

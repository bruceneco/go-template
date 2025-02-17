package grpc

import (
	"github.com/bruceneco/go-template/internal/adapters/grpc/servers"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func() *validator.Validate {
		return validator.New(validator.WithRequiredStructEnabled())
	}),
	fx.Provide(
		fx.Annotate(servers.NewUserServer, fx.As(new(Server)), fx.ResultTags(`group:"servers"`)),
	),
	fx.Provide(NewGRPCServer),
)

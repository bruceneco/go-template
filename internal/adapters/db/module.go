package db

import (
	"github.com/bruceneco/go-template/internal/adapters/db/postgres"
	"github.com/bruceneco/go-template/internal/domain/user"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		postgres.NewConnection,
		fx.Annotate(
			postgres.NewUserRepository,
			fx.As(new(user.Repository)),
		),
	),
)

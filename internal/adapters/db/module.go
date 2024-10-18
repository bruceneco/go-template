package db

import (
	"go-template/internal/adapters/db/postgres"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(postgres.NewConnection),
)

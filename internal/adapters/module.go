package adapters

import (
	"go-template/internal/adapters/amqp"
	"go-template/internal/adapters/db/postgres"
	"go-template/internal/adapters/http"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		postgres.NewConnection,
		amqp.NewConnection,
	),
	http.Module,
)

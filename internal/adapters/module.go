package adapters

import (
	"go-template/internal/adapters/amqp"
	"go-template/internal/adapters/amqp/consumers"
	"go-template/internal/adapters/amqp/publishers"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		amqp.NewConnection,
		publishers.NewExample,
	),
	fx.Invoke(consumers.StartExample),
)

package publishers

import (
	"context"
	"go-template/internal/adapters/amqp"

	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/fx"
)

type Example struct {
	p *rabbitmq.Publisher
}

func NewExample(conn *amqp.Connection, lc fx.Lifecycle) *Example {
	publisher, err := conn.MakePublisher(&amqp.ExchangeOpts{
		Name: amqp.DefaultExchangeName,
		Kind: amqp.DefaultExchangeType,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("cannot initialize example amqp publisher")
	}
	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			publisher.Close()
			return nil
		},
	})

	return &Example{p: publisher}
}

func (e *Example) Publish(ctx context.Context) error {
	return e.p.PublishWithContext(ctx, []byte("example"), []string{"example_queue"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
}

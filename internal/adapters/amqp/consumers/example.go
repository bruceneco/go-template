package consumers

import (
	"context"
	"go-template/internal/adapters/amqp"

	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/fx"
)

func StartExample(lc fx.Lifecycle, conn *amqp.Connection) {
	const numConcurrency = 5
	consumer, err := conn.MakeConsumer("example_consumer", &amqp.ConsumerOpts{
		QueueName:   "example_queue",
		Concurrency: numConcurrency,
	}, &amqp.ExchangeOpts{
		Name: amqp.DefaultExchangeName,
		Kind: amqp.DefaultExchangeType,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("cannot initialize example amqp consumer")
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			consumer.CloseWithContext(ctx)
			return nil
		},
	})
	go func() {
		err = consumer.Run(Handle)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot run example amqp consumer")
		}
	}()
}

func Handle(d rabbitmq.Delivery) rabbitmq.Action {
	log.Info().Str("content", string(d.Body)).Msg("got example message from example queue")
	return rabbitmq.Ack
}

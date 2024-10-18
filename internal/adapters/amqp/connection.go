package amqp

import (
	"context"
	"go-template/config"

	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/fx"
)

type Connection struct {
	conn *rabbitmq.Conn
}

func NewConnection(lc fx.Lifecycle, cfg *config.EnvConfig) *Connection {
	conn, err := rabbitmq.NewConn(
		cfg.AMQPHost,
		rabbitmq.WithConnectionOptionsLogger(NewLoggerAdapter(&log.Logger)),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot initialize amqp connection")
	}
	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			if conn == nil {
				return nil
			}
			return conn.Close()
		},
	})
	return &Connection{conn: conn}
}

type ExchangeOpts struct {
	Name ExchangeName
	Kind ExchangeType
}

func (c *Connection) MakePublisher(opts *ExchangeOpts) (*rabbitmq.Publisher, error) {
	p, err := rabbitmq.NewPublisher(
		c.conn,
		rabbitmq.WithPublisherOptionsLogger(NewLoggerAdapter(&log.Logger)),
		rabbitmq.WithPublisherOptionsExchangeName(opts.Name.String()),
		rabbitmq.WithPublisherOptionsExchangeDurable,
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeKind(opts.Kind.String()),
	)
	if err != nil {
		return nil, err
	}
	p.NotifyReturn(func(r rabbitmq.Return) {
		log.Info().Interface("content", r).Msg("message returned from exchange")
	})
	p.NotifyPublish(func(c rabbitmq.Confirmation) {
		log.Info().Interface("confirmation", c).Msg("message confirmed")
	})
	return p, nil
}

type ConsumerOpts struct {
	QueueName   string
	Exclusive   bool
	NoWait      bool
	Concurrency int
}

func (c *Connection) MakeConsumer(name string, cOpts *ConsumerOpts, exOpts *ExchangeOpts) (*rabbitmq.Consumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		c.conn,
		name,
		rabbitmq.WithConsumerOptionsLogger(NewLoggerAdapter(&log.Logger)),
		rabbitmq.WithConsumerOptionsExchangeName(exOpts.Name.String()),
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeKind(exOpts.Kind.String()),
		func(options *rabbitmq.ConsumerOptions) {
			options.QueueOptions = rabbitmq.QueueOptions{
				Name:       cOpts.QueueName,
				Durable:    true,
				AutoDelete: false,
				Exclusive:  cOpts.Exclusive,
				Passive:    false,
				NoWait:     cOpts.NoWait,
				Declare:    true,
			}
			options.CloseGracefully = true
			if cOpts.Concurrency != 0 {
				options.Concurrency = cOpts.Concurrency
			}
		},
	)
	return consumer, err
}

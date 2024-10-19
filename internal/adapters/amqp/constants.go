package amqp

import "github.com/rabbitmq/amqp091-go"

type ExchangeName string

const (
	DefaultExchangeName ExchangeName = "events"
)

type ExchangeType string

const (
	DefaultExchangeType ExchangeType = amqp091.ExchangeHeaders
	ExchangeTypeHeaders
)

func (n *ExchangeName) String() string {
	return string(*n)
}

func (t *ExchangeType) String() string {
	return string(*t)
}

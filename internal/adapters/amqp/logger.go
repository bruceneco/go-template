package amqp

import (
	"github.com/rs/zerolog"
	"github.com/wagslane/go-rabbitmq"
)

var _ rabbitmq.Logger = (*LoggerAdapter)(nil)

type LoggerAdapter struct {
	log *zerolog.Logger
}

func NewLoggerAdapter(log *zerolog.Logger) *LoggerAdapter {
	return &LoggerAdapter{log: log}
}
func (l *LoggerAdapter) Fatalf(s string, i ...interface{}) {
	l.log.Fatal().Msgf(s, i...)
}

func (l *LoggerAdapter) Errorf(s string, i ...interface{}) {
	l.log.Error().Msgf(s, i...)
}

func (l *LoggerAdapter) Warnf(s string, i ...interface{}) {
	l.log.Warn().Msgf(s, i...)
}

func (l *LoggerAdapter) Infof(_ string, _ ...interface{}) {}

func (l *LoggerAdapter) Debugf(s string, i ...interface{}) {
	l.log.Debug().Msgf(s, i...)
}

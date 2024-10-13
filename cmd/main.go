package main

import (
	"context"
	"time"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Invoke(NewExample),
	).Run()
}

func NewExample(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return nil
		},
		OnStop: func(_ context.Context) error {
			gracefulDelay := 3
			time.Sleep(time.Duration(gracefulDelay) * time.Second)
			return nil
		},
	})
}

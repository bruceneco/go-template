package http

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHTTPServer),
	fx.Invoke(NewHealthCheckController),
)

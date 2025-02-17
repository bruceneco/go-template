package adapters

import (
	"github.com/bruceneco/go-template/internal/adapters/amqp"
	"github.com/bruceneco/go-template/internal/adapters/db"
	"github.com/bruceneco/go-template/internal/adapters/grpc"
	"github.com/bruceneco/go-template/internal/adapters/http"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		amqp.NewConnection,
	),
	db.Module,
	grpc.Module,
	http.Module,
)

package adapters

import (
	"go-template/internal/adapters/db"
	"go.uber.org/fx"
)

var Module = fx.Options(
	db.Module,
)

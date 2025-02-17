package domain

import (
	"go-template/internal/domain/user"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	user.NewService,
)

package models

import "errors"

var (
	ErrNilModel          = errors.New("model is nil")
	ErrNilEntity         = errors.New("entity is nil")
	ErrCantStorePassword = errors.New("can't store password")
	ErrNilID             = errors.New("id is nil")
)

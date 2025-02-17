package user

import "errors"

var (
	ErrInvalidID          = errors.New("invalid id")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNilEntity          = errors.New("entity is nil")
)

package models

import (
	"fmt"
	"github.com/bruceneco/go-template/internal/domain/user"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type User struct {
	Base
	Name     string
	Email    string
	Password string
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u == nil {
		log.Error().Stack().Msg("model is nil")
		return ErrNilModel
	}

	return u.NewID()
}
func (u *User) FromEntity(entity *user.Entity) (*User, error) {
	*u = User{
		Name:  entity.Name,
		Email: entity.Email,
	}
	u.ID = entity.ID
	hashedPass, err := entity.PasswordHash()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantStorePassword, err)
	}
	u.Password = hashedPass
	return u, nil
}
func (u *User) ToEntity() user.Entity {
	return user.Entity{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

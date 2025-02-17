package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Entity struct {
	ID       *uuid.UUID
	Name     string
	Email    string
	Password string
}

func (e *Entity) PasswordHash() (string, error) {
	if _, err := bcrypt.Cost([]byte(e.Password)); err == nil {
		return e.Password, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (e *Entity) ComparePassword(p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(p))
	return err == nil
}

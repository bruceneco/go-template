package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        *uuid.UUID `gorm:"type:UUID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *Base) NewID() error {
	if b == nil {
		return fmt.Errorf("base is nil: %w", ErrNilModel)
	}
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	b.ID = &id
	return nil
}

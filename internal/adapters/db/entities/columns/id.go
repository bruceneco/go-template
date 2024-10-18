package columns

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ID struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (id *ID) BeforeCreate(_ *gorm.DB) error {
	generated, err := uuid.NewV6()
	if err != nil {
		return err
	}
	id.ID = generated
	return nil
}

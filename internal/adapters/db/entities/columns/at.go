package columns

import (
	"time"

	"gorm.io/gorm"
)

type CreatedAt struct {
	CreatedAt time.Time `gorm:"created_at"`
}

type UpdatedAt struct {
	UpdatedAt time.Time `gorm:"created_at"`
}

type DeletedAt struct {
	DeletedAt gorm.DeletedAt `gorm:"deleted_at"`
}

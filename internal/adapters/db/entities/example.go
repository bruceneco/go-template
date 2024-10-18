package entities

import (
	"go-template/internal/adapters/db/entities/columns"
)

type Example struct {
	columns.ID
	columns.CreatedAt
	columns.UpdatedAt
	columns.DeletedAt
}

func (e *Example) TableName() string {
	return "examples"
}

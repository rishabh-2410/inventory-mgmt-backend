package models

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	ID         uuid.UUID
	BusinessID uuid.UUID
	Name       string
	Address    *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

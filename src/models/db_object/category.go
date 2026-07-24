package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID
	Name        string
	BusinessID  uuid.UUID
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

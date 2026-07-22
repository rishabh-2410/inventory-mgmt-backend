package models

import (
	"time"

	"github.com/google/uuid"
)

type Businesss struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Phone     *string
	Address   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

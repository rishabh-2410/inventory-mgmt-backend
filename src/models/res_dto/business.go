package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateBusinessResponseDTO struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Phone     *string
	Address   *string
	CreatedAt time.Time
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateBusinessResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone"`
	Address   *string   `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

package models

import (
	"github.com/google/uuid"
)

type AddWarehouseResponseDTO struct {
	ID         uuid.UUID `json:"id"`
	BusinessID uuid.UUID `json:"business_id"`
	Name       string    `json:"name"`
	Address    *string   `json:"address"`
}

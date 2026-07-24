package models

import "github.com/google/uuid"

type AddCategoryResponseDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	BusinessID  uuid.UUID `json:"business_id"`
	Description *string   `json:"description"`
}

package models

import "github.com/google/uuid"

type CreateBusinessRequestDTO struct {
	UserID  uuid.UUID `json:"user_id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Phone   *string   `json:"phone"`
	Address *string   `json:"address"`
}

package models

import (
	"time"

	enums "inventory-backend/src/models/enums"

	"github.com/google/uuid"
)

type LoginUserResponseDTO struct {
	ID                   uuid.UUID       `json:"id"`
	Name                 string          `json:"name"`
	Email                string          `json:"email"`
	BusinessID           *uuid.UUID      `json:"business_id"`
	Role                 *enums.UserRole `json:"role"`
	IsActive             bool            `json:"is_active"`
	AccessToken          string          `json:"access_token"`
	AccessTokenExpiresAt time.Time       `json:"access_token_expires_at"`
	RefreshToken         string          `json:"refresh_token"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type RegisterUserResponseDTO struct {
	ID uuid.UUID `json:"id"`
}

type CreateUserResponseDTO struct {
	ID         uuid.UUID       `json:"id"`
	Email      string          `json:"email"`
	Name       string          `json:"name"`
	BusinessID *uuid.UUID      `json:"business_id"`
	Role       *enums.UserRole `json:"role"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

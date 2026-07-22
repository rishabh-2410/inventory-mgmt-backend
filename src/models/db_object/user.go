package models

import (
	enums "inventory-backend/src/models/enums"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	BusinessID   *uuid.UUID
	PasswordHash string
	Role         *enums.UserRole
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

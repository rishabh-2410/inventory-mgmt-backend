package auth

import (
	enums "inventory-backend/src/models/enums"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Role       enums.UserRole `json:"role"`
	BusinessID uuid.UUID      `json:"business_id"`
	jwt.RegisteredClaims
}

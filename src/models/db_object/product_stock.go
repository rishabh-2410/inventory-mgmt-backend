package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductStock struct {
	ID           uuid.UUID
	BusinessID   uuid.UUID
	ProductID    uuid.UUID
	WarehouseID  uuid.UUID
	CurrentStock float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

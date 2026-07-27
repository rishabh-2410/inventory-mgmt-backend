package models

import (
	"time"

	"github.com/google/uuid"
)

type StockMovement struct {
	ID               uuid.UUID
	ProductStockID   uuid.UUID
	UserID           uuid.UUID
	ProductID        uuid.UUID
	WarehouseID      uuid.UUID
	BusinessID       uuid.UUID
	MovementType     string
	MovementQuantity float64
	Reason           string
	CreatedAt        time.Time
}

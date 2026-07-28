package models

import (
	"time"

	"github.com/google/uuid"
)

type StockMovementResponseDTO struct {
	ProductStockID   uuid.UUID
	ProductID        uuid.UUID
	WarehouseID      uuid.UUID
	CurrentStock     float64
	MovementType     string
	MovementQuantity float64
	Reason           string
	CreatedAt        time.Time
}

type StockMovementHistoryResponseDTO struct {
	UserID           uuid.UUID
	WarehouseID      uuid.UUID
	ProductID        uuid.UUID
	MovementType     string
	MovementQuantity float64
	Reason           string
	CreatedAt        time.Time
}

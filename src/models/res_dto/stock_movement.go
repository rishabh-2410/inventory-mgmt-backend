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

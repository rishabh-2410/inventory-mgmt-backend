package models

import (
	"github.com/google/uuid"
)

type StockMovementRequestDTO struct {
	ProductID        uuid.UUID
	WarehouseID      uuid.UUID
	MovementType     string
	MovementQuantity float64
	Reason           string
}

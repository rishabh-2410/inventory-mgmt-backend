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
	UserID           uuid.UUID `json:"user_id"`
	ProductID        uuid.UUID `json:"product_id"`
	MovementType     string    `json:"movement_type"`
	MovementQuantity float64   `json:"movement_quantity"`
	Reason           string    `json:"reason"`
	CreatedAt        time.Time `json:"created_at"`
	ProductName      string    `json:"product_name"`
	Sku              string    `json:"sku"`
	WarehouseName    string    `json:"warehouse_name"`
}

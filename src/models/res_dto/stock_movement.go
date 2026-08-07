package models

import (
	"time"

	"github.com/google/uuid"
)

type StockMovementResponseDTO struct {
	ProductStockID   uuid.UUID `json:"product_stock_id"`
	ProductID        uuid.UUID `json:"product_id"`
	WarehouseID      uuid.UUID `json:"warehouse_id"`
	CurrentStock     float64   `json:"current_stock"`
	MovementType     string    `json:"movement_type"`
	MovementQuantity float64   `json:"movement_quantity"`
	Reason           string    `json:"reason"`
	CreatedAt        time.Time `json:"created_at"`
}

type StockMovementHistory struct {
	ID               uuid.UUID `json:"id"`
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

type StockMovementHistoryResponseDTO struct {
	Data    []*StockMovementHistory `json:"data"`
	Total   int                     `json:"total"`
	Page    int                     `json:"page"`
	Limit   int                     `json:"limit"`
	HasNext bool                    `json:"has_next"`
}

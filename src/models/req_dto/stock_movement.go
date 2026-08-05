package models

import (
	"time"

	"github.com/google/uuid"
)

type StockMovementRequestDTO struct {
	ProductID        uuid.UUID
	WarehouseID      uuid.UUID
	MovementType     string
	MovementQuantity float64
	Reason           string
}

type GetStockMovementHistoryRequestDTO struct {
	Limit        *int       `json:"limit"`
	Page         *int       `json:"page"`
	MovementType *string    `json:"movement_type"`
	WarehouseID  *uuid.UUID `json:"warehouse_id"`
	ProductID    *uuid.UUID `json:"product_id"`
	Search       *string    `json:"search"`
	From         *time.Time `json:"from"`
	To           *time.Time `json:"to"`
	SortBy       *string    `json:"sort_by"`
}

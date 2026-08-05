package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductStockResponseDTO struct {
	ID            uuid.UUID `json:"id"`
	ProductName   string    `json:"product_name"`
	Sku           string    `json:"sku"`
	CategoryName  string    `json:"category_name"`
	WarehouseName string    `json:"warehouse_name"`
	CurrentStock  float64   `json:"current_stock"`
	UpdatedAt     time.Time `json:"updated_at"`
}

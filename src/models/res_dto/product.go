package models

import (
	"time"

	"github.com/google/uuid"
)

type AddProductResponseDTO struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	CategoryID   uuid.UUID `json:"category_id"`
	BusinessID   uuid.UUID `json:"business_id"`
	SKU          string    `json:"sku"`
	SellingPrice float64   `json:"selling_price"`
	CostPrice    float64   `json:"cost_price"`
	ImageURL     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductWithStock struct {
	AddProductResponseDTO
	TotalStock float64 `json:"total_stock"`
}

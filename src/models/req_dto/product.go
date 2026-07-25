package models

import "github.com/google/uuid"

type AddProductRequestDTO struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	CategoryID   uuid.UUID `json:"category_id"`
	BusinessID   uuid.UUID `json:"business_id"`
	SKU          string    `json:"sku"`
	SellingPrice float64   `json:"selling_price"`
	CostPrice    float64   `json:"cost_price"`
	ImageURL     string    `json:"image_url"`
}

type EditProductRequestDTO struct {
	Name         *string    `json:"name"`
	CategoryID   *uuid.UUID `json:"category_id"`
	SKU          *string    `json:"sku"`
	SellingPrice *float64   `json:"selling_price"`
	CostPrice    *float64   `json:"cost_price"`
	ImageURL     *string    `json:"image_url"`
}

package models

import "github.com/google/uuid"

type GetInventoryListRequestDTO struct {
	Search      *string    `json:"search"`
	CategoryID  *uuid.UUID `json:"category_id"`
	WarehouseID *uuid.UUID `json:"warehouse_id"`
	ProductID   *uuid.UUID `json:"product_id"`
	Limit       *int       `json:"limit"`
	Page        *int       `json:"page"`
}

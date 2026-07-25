package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID
	Name         string
	CategoryID   uuid.UUID
	BusinessID   uuid.UUID
	SKU          string
	SellingPrice float64
	CostPrice    float64
	ImageURL     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

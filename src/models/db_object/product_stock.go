package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductStock struct {
	ID            uuid.UUID
	ProductName   string
	Sku           string
	CategoryName  string
	WarehouseName string
	CurrentStock  float64
	UpdatedAt     time.Time
}

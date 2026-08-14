package models

type AddWarehouseRequestDTO struct {
	Name    string  `json:"name" validate:"required"`
	Address *string `json:"address" validate:"omitempty,gt=0"`
}

type EditWarehouseRequestDTO struct {
	Name    *string `json:"name" validate:"omitempty"`
	Address *string `json:"address" validate:"omitempty"`
}

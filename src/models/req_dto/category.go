package models

type AddCategoryRequestDTO struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description" validate:"omitempty,gt=0"`
}

type UpdateCategoryRequestDTO struct {
	Name        *string `json:"name" validate:"omitempty,gt=0"`
	Description *string `json:"description" validate:"omitempty,gt=0"`
}

package models

type CreateBusinessRequestDTO struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`
}

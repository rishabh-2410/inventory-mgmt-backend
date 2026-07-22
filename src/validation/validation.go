package validation

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func InitializeValidator() {
	// Here we're using a single instance of validate object (global)
	// It caches struct info. So any validation and parsing of struct is only done once per struct.
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

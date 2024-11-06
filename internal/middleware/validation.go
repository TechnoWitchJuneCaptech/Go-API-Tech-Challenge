package validation

import "github.com/go-playground/validator"

// ValidateType is a custom validation function for ../models.Person objects ensuring the Type field is either "professor" or "student".
func ValidateType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == "professor" || value == "student"
}

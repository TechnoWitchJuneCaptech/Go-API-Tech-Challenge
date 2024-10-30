package validation

import "github.com/go-playground/validator"

func ValidateType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == "professor" || value == "student"
}

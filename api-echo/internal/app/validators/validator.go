package validators

import (
	"fmt"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate() map[string][]string
}

func ValidateRequest(req Validator) map[string][]string {
	return req.Validate()
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func StructToErrors(err error) map[string][]string {
	errors := make(map[string][]string)
	if err == nil {
		return errors
	}

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		message := fmt.Sprintf("%s validation failed on %s", field, tag)
		errors[field] = append(errors[field], message)
	}
	return errors
}

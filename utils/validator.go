package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func Validate(value any) error {
	err := requestValidator.Struct(value)
	if err == nil {
		return nil
	}
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	messages := make([]string, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		field := strings.ReplaceAll(fieldError.Field(), "_", " ")
		switch fieldError.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", field))
		case "gte":
			messages = append(messages, fmt.Sprintf("%s must be at least %s", field, fieldError.Param()))
		case "gt":
			messages = append(messages, fmt.Sprintf("%s must be greater than %s", field, fieldError.Param()))
		case "email":
			messages = append(messages, fmt.Sprintf("%s must be a valid email address", field))
		default:
			messages = append(messages, fmt.Sprintf("%s is invalid", field))
		}
	}
	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

package validations

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
)

// This is the validator instance
// for more information see: https://github.com/go-playground/validator
var validate = validator.New()

type ValidationError struct {
	FailedField string
	Tag         string
	Param       string
	Value       any
}

func Validate(data any) error {
	var validationErrors []ValidationError

	// Perform validation using the validator instance
	if errs := validate.Struct(data); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Value(),
				Param:       err.Param(),
			})
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	// Convert validation errors into a single error message
	return errors.New(strings.Join(translateErrors(validationErrors), "; "))
}

func translateErrors(errs []ValidationError) []string {
	var errorMessages []string
	for _, err := range errs {
		switch err.Tag {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("[%s] is required", err.FailedField))
		case "min":
			errorMessages = append(errorMessages, fmt.Sprintf("[%s] must be at least %s characters long", err.FailedField, err.Param))
		case "max":
			errorMessages = append(errorMessages, fmt.Sprintf("[%s] must be at most %s characters long", err.FailedField, err.Param))
		case "email":
			errorMessages = append(errorMessages, fmt.Sprintf("[%s] must be a valid email address", err.FailedField))
		case "gte":
			errorMessages = append(errorMessages, fmt.Sprintf("[%s] must be greater than or equal to %s", err.FailedField, err.Param))
		case "lte":
			errorMessages = append(errorMessages, fmt.Sprintf("[%s] must be less than or equal to %s", err.FailedField, err.Param))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("[%s] is not valid", err.FailedField))
		}
	}
	return errorMessages
}

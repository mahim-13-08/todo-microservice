package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate is the singleton validator instance.
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

// ValidationErrors formats validation errors into a human-readable string.
func ValidationErrors(err error) string {
	var errs []string
	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, formatError(e))
	}
	return strings.Join(errs, "; ")
}

func formatError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("field '%s' is required", strings.ToLower(e.Field()))
	case "min":
		return fmt.Sprintf("field '%s' must be at least %s characters", strings.ToLower(e.Field()), e.Param())
	case "max":
		return fmt.Sprintf("field '%s' must be at most %s characters", strings.ToLower(e.Field()), e.Param())
	case "oneof":
		return fmt.Sprintf("field '%s' must be one of [%s]", strings.ToLower(e.Field()), e.Param())
	default:
		return fmt.Sprintf("field '%s' failed validation on '%s'", strings.ToLower(e.Field()), e.Tag())
	}
}

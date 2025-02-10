package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/urlpick/urlpick-api/internal/pkg/utils/errors"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Struct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fieldErrors := make(map[string]string)
			for _, e := range validationErrors {
				fieldErrors[e.Field()] = getErrorMsg(e)
			}
			return errors.BadRequest("Validation failed").WithDetails(fieldErrors)
		}
		return errors.BadRequest("Invalid input")
	}
	return nil
}

func getErrorMsg(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "url":
		return "Invalid URL format"
	default:
		return "Invalid value"
	}
}

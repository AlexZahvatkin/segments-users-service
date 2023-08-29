package validator_utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationError(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is too short", err.Field()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is too long", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return strings.Join(errMsgs, ", ")
}

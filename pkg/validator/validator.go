package validator

import (
	"github.com/go-playground/validator/v10"
)

func ParseValidationErrors(err error) []string {
	var errMsgs []string
	if vErrs, ok := err.(validator.ValidationErrors); ok {
		for _, vErr := range vErrs {
			switch vErr.Tag() {
			case "required":
				errMsgs = append(errMsgs, vErr.Field()+": is required")
			case "email":
				errMsgs = append(errMsgs, vErr.Field()+": invalid email format")
			case "uuid":
				errMsgs = append(errMsgs, vErr.Field()+": invalid uuid format")
			case "min":
				errMsgs = append(errMsgs, vErr.Field()+": is too short")
			case "max":
				errMsgs = append(errMsgs, vErr.Field()+": is too long")
			case "startswith":
				errMsgs = append(errMsgs, vErr.Field()+": must start with "+vErr.Param())
			case "len":
				errMsgs = append(errMsgs, vErr.Field()+": must be exactly "+vErr.Param()+" characters")
			default:
				errMsgs = append(errMsgs, vErr.Field()+": invalid value")
			}
		}
	}
	return errMsgs
}

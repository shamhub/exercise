package errorlib

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateQueryParams(data interface{}) error {
	err := validate.Struct(data)

	if err != nil {
		return sendErrorList(err)
	}
	return nil
}

func sendErrorList(err error) error {

	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorList []*ResponseError

		for _, fieldErr := range validationErrors {
			errorMessage := fmt.Sprintf("Field: %s, Tag: %s, Value: %v\n", fieldErr.Field(), fieldErr.Tag(), fieldErr.Value())
			fieldError := NewResponseError(400, errorMessage)
			errorList = append(errorList, fieldError)
		}
		return NewMultiError(400, errorList)
	}
	return NewResponseError(400, err.Error())
}

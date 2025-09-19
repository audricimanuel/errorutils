package errorutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strings"
	"time"
)

// GetValidatorController return validator controller
func GetValidatorController() *validator.Validate {
	return validator.New()
}

// getErrorMessage to define the error when validating struct field with tag 'validate'.
//
//	Usage example:
//		errorMessage := getErrorMessage(err, "required")
//		return errorMessage
//
//	Define new error:
//		case "error_something":
//			return "this error occurs because of this case is not fulfilled"
func getErrorMessage(err validator.FieldError, jsonField string) string {
	if jsonField == "" {
		return err.Tag()
	}

	switch err.Tag() {
	case "required", "required_if":
		return fmt.Sprintf("%s is required", jsonField)
	case "min":
		return fmt.Sprintf("%s must be at least %s", jsonField, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", jsonField, err.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", jsonField)
	case "oneof":
		choices := strings.Split(err.Param(), " ")
		choicesStr := ""
		for i, v := range choices {
			choicesStr += fmt.Sprintf(`%v`, v)
			if i != len(choices)-1 {
				choicesStr += ", "
			}
		}
		return fmt.Sprintf("%s valid choices are: %s", jsonField, choicesStr)
	case "datetime_format":
		return fmt.Sprintf("%s datetime format is YYYY-MM-DD hh:mm:ss", jsonField)
	case "date_format":
		return fmt.Sprintf("%s date format is YYYY-MM-DD", jsonField)
	case "gt":
		minimumLength := err.Param()
		if minimumLength == "0" {
			return fmt.Sprintf("%s can't be empty", jsonField)
		}
		return fmt.Sprintf("Minimum length of %s is %s", jsonField, minimumLength)
	default:
		return fmt.Sprintf("Validation error on field %s", jsonField)
	}
}

// ValidatePayload to validate payload when using tag 'validate' (returning HttpError)
func ValidatePayload(request *http.Request, s interface{}) HttpError {
	err := json.NewDecoder(request.Body).Decode(s)
	if err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			errorType := err.(*json.UnmarshalTypeError)
			errorFormat := ErrorBadRequest.CustomMessage(fmt.Sprintf("invalid type of %s (expected: %s, got: %s)", errorType.Field, errorType.Type, errorType.Value))
			return errorFormat
		default:
			return ErrorBadRequest.CustomMessage(fmt.Sprintf("payload error: %s", err.Error()))
		}
	}

	if err := ValidateStruct(s); err != nil {
		return ErrorBadRequest.CustomMessage(err.Error())
	}

	return nil
}

// ValidateStruct to validate struct using Go Validator
func ValidateStruct(structObj interface{}) error {
	validatorObj := GetValidatorController()

	// add custom validator to validate field with datetime format (YYYY-MM-DD hh:mm)
	validatorObj.RegisterValidation("datetime_format", func(fl validator.FieldLevel) bool {
		datetimeStr := fl.Field().String()
		_, err := time.Parse(FORMAT_DATETIME_DEFAULT, datetimeStr)
		return err == nil
	})

	// add custom validator to validate field with date format (YYYY-MM-DD)
	validatorObj.RegisterValidation("date_format", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		_, err := time.Parse(FORMAT_DATE_DEFAULT, dateStr)
		return err == nil
	})

	if err := validatorObj.Struct(structObj); err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorType := err.(validator.ValidationErrors)
			for i := 0; i < len(errorType); {
				errorField := errorType[i]
				jsonField := GetJsonTagInStruct(errorField.Field(), structObj)
				errorMessage := getErrorMessage(errorField, jsonField)
				return errors.New(errorMessage)
			}
		case *validator.InvalidValidationError:
			errorType := err.(*validator.InvalidValidationError)
			log.Println(fmt.Sprintf(`[ERROR] validator.InvalidValidationError: {"type": "%v", "key": "%v", "name": "%v"}`, errorType.Type, errorType.Type.Key(), errorType.Type.Name()))
			return errors.New(fmt.Sprintf("payload error: %s", errorType))
		default:
			return errors.New(fmt.Sprintf("payload error: %s", err.Error()))
		}
	}
	return nil
}

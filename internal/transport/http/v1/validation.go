package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"go-rest-api-example/internal/pkg/logger"
)

// validator instance for the v1 package
var validate = validator.New()

// ValidationError represents a validation error response
type ValidationError struct {
	Message string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// DecodeAndValidateJSON decodes JSON from request body and validates it
func DecodeAndValidateJSON(r *http.Request, target interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return &ValidationError{
			Message: "invalid JSON format",
		}
	}

	if err := validate.Struct(target); err != nil {
		details := make(map[string]string)
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				field := fieldError.Field()
				tag := fieldError.Tag()
				details[field] = getValidationMessage(field, tag, fieldError.Param())
			}
		}

		return &ValidationError{
			Message: "validation failed",
			Details: details,
		}
	}

	return nil
}

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, r *http.Request, err *ValidationError) {
	log := logger.FromContext(r.Context())
	if log != nil {
		log.Warn("validation error", zap.Any("details", err.Details))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(err)
}

func getValidationMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "min":
		if param != "" {
			return field + " must be at least " + param + " characters long"
		}
		return field + " is too short"
	case "max":
		if param != "" {
			return field + " must be no more than " + param + " characters long"
		}
		return field + " is too long"
	case "email":
		return field + " must be a valid email address"
	case "url":
		return field + " must be a valid URL"
	default:
		return field + " is invalid"
	}
}

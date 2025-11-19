package middleware

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/NoroSaroyan/go-rest-api-example/internal/pkg/logger"
)

const (
	// jsonTagParts defines how many parts we expect when splitting JSON tags
	jsonTagParts = 2
)

// Validator wraps the go-playground/validator for dependency injection
type Validator interface {
	Validate(s interface{}) error
}

type customValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() Validator {
	v := validator.New()

	// Register custom tag name function to use JSON tag names in error messages
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", jsonTagParts)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &customValidator{validator: v}
}

func (cv *customValidator) Validate(s interface{}) error {
	return cv.validator.Struct(s)
}

// ValidationErrorResponse represents validation error details
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details"`
}

// ValidateJSON middleware validates JSON request bodies against struct validation tags
func ValidateJSON(v Validator, target interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())

			// Create a new instance of the target type
			targetType := reflect.TypeOf(target)
			if targetType.Kind() == reflect.Ptr {
				targetType = targetType.Elem()
			}
			targetValue := reflect.New(targetType).Interface()

			// Decode JSON body
			if err := json.NewDecoder(r.Body).Decode(targetValue); err != nil {
				if log != nil {
					log.Warn("invalid JSON body", zap.Error(err))
				}
				writeValidationError(w, "invalid JSON format", nil)
				return
			}

			// Validate the decoded struct
			if err := v.Validate(targetValue); err != nil {
				if log != nil {
					log.Warn("validation failed", zap.Error(err))
				}

				details := make(map[string]string)
				if validationErrors, ok := err.(validator.ValidationErrors); ok {
					for _, fieldError := range validationErrors {
						field := fieldError.Field()
						tag := fieldError.Tag()
						details[field] = getValidationMessage(field, tag, fieldError.Param())
					}
				}

				writeValidationError(w, "validation failed", details)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func writeValidationError(w http.ResponseWriter, message string, details map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := ValidationErrorResponse{
		Error:   message,
		Details: details,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If we can't encode the validation error response, there's not much we can do
		// except log the error - the status code has already been set
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
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

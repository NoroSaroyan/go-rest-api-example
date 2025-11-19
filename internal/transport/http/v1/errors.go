package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"go-rest-api-example/internal/domain"
	"go-rest-api-example/internal/pkg/logger"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
}

// AppError represents an internal application error with context
type AppError struct {
	Err        error
	Message    string
	Code       string
	HTTPStatus int
	Context    map[string]interface{}
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(err error, message, code string, httpStatus int) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		Code:       code,
		HTTPStatus: httpStatus,
		Context:    make(map[string]interface{}),
	}
}

// WithContext adds context to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	e.Context[key] = value
	return e
}

// WriteJSONSafe writes JSON response and handles encoding errors gracefully
func WriteJSONSafe(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log := logger.FromContext(r.Context())
		if log != nil {
			log.Error("failed to encode JSON response",
				zap.Error(err),
				zap.Any("response", v),
				zap.Int("status_code", code),
			)
		}

		// If JSON encoding fails, write a minimal error response
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(`{"error":"internal server error","code":"ENCODING_ERROR"}`)); writeErr != nil {
			if log != nil {
				log.Error("failed to write error response", zap.Error(writeErr))
			}
		}
	}
}

// WriteError handles error responses with proper logging and security
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	log := logger.FromContext(r.Context())
	traceID := getTraceID(r)

	var appErr *AppError
	var response ErrorResponse

	if errors.As(err, &appErr) {
		// Handle application errors
		response = ErrorResponse{
			Error:   appErr.Error(),
			Code:    appErr.Code,
			TraceID: traceID,
		}

		if log != nil {
			fields := []zap.Field{
				zap.Error(appErr.Err),
				zap.String("app_error_code", appErr.Code),
				zap.String("trace_id", traceID),
				zap.Int("http_status", appErr.HTTPStatus),
			}

			// Add context fields
			for key, value := range appErr.Context {
				fields = append(fields, zap.Any(key, value))
			}

			log.Error("application error", fields...)
		}

		WriteJSONSafe(w, r, appErr.HTTPStatus, response)
		return
	}

	// Handle domain errors
	switch {
	case errors.Is(err, domain.ErrTodoNotFound):
		response = ErrorResponse{
			Error:   "todo not found",
			Code:    "TODO_NOT_FOUND",
			TraceID: traceID,
		}

		if log != nil {
			log.Warn("todo not found",
				zap.Error(err),
				zap.String("trace_id", traceID),
			)
		}

		WriteJSONSafe(w, r, http.StatusNotFound, response)

	case errors.Is(err, domain.ErrInvalidTitle):
		response = ErrorResponse{
			Error:   "title cannot be empty",
			Code:    "INVALID_TITLE",
			TraceID: traceID,
		}

		if log != nil {
			log.Warn("invalid title provided",
				zap.Error(err),
				zap.String("trace_id", traceID),
			)
		}

		WriteJSONSafe(w, r, http.StatusBadRequest, response)

	default:
		// Handle unexpected errors - log full details but return generic message
		response = ErrorResponse{
			Error:   "internal server error",
			Code:    "INTERNAL_ERROR",
			TraceID: traceID,
		}

		if log != nil {
			log.Error("unexpected error occurred",
				zap.Error(err),
				zap.String("trace_id", traceID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)
		}

		WriteJSONSafe(w, r, http.StatusInternalServerError, response)
	}
}

// getTraceID extracts trace ID from request context or generates a fallback
func getTraceID(r *http.Request) string {
	if traceID := r.Header.Get("X-Trace-Id"); traceID != "" {
		return traceID
	}

	// Try to get request ID if available
	if reqID, ok := r.Context().Value("request_id").(string); ok && reqID != "" {
		return reqID
	}

	return "unknown"
}

// Database error helpers
func NewDatabaseError(err error, operation string) *AppError {
	return NewAppError(err, "database operation failed", "DATABASE_ERROR", http.StatusInternalServerError).
		WithContext("operation", operation)
}

// Validation error helpers
func NewValidationError(message string) *AppError {
	return NewAppError(nil, message, "VALIDATION_ERROR", http.StatusBadRequest)
}

// Not found error helpers
func NewNotFoundError(resource string) *AppError {
	return NewAppError(nil, resource+" not found", "NOT_FOUND", http.StatusNotFound).
		WithContext("resource", resource)
}

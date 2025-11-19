package v1

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeAndValidateJSON(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		target      interface{}
		wantErr     bool
		expectCode  string
		description string
	}{
		{
			name:        "valid request",
			body:        `{"title": "Test Todo"}`,
			target:      &CreateTodoRequest{},
			wantErr:     false,
			description: "should successfully decode and validate valid request",
		},
		{
			name:        "missing required field",
			body:        `{}`,
			target:      &CreateTodoRequest{},
			wantErr:     true,
			description: "should fail validation when required field is missing",
		},
		{
			name:        "empty title",
			body:        `{"title": ""}`,
			target:      &CreateTodoRequest{},
			wantErr:     true,
			description: "should fail validation when title is empty",
		},
		{
			name:        "title too long",
			body:        `{"title": "` + generateLongString(300) + `"}`,
			target:      &CreateTodoRequest{},
			wantErr:     true,
			description: "should fail validation when title exceeds maximum length",
		},
		{
			name:        "invalid JSON",
			body:        `{"title": }`,
			target:      &CreateTodoRequest{},
			wantErr:     true,
			description: "should fail with invalid JSON format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			err := DecodeAndValidateJSON(req, tt.target)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DecodeAndValidateJSON() expected error but got nil - %s", tt.description)
				}

				if validationErr, ok := err.(*ValidationError); ok {
					if validationErr.Message == "" {
						t.Error("ValidationError should have a non-empty error message")
					}
				}
				return
			}

			if err != nil {
				t.Errorf("DecodeAndValidateJSON() unexpected error = %v - %s", err, tt.description)
				return
			}

			// For successful cases, verify the data was decoded
			if req, ok := tt.target.(*CreateTodoRequest); ok {
				if req.Title != "Test Todo" {
					t.Errorf("DecodeAndValidateJSON() title = %v, want 'Test Todo'", req.Title)
				}
			}
		})
	}
}

func TestWriteValidationError(t *testing.T) {
	tests := []struct {
		name        string
		err         *ValidationError
		wantStatus  int
		description string
	}{
		{
			name: "simple validation error",
			err: &ValidationError{
				Message: "validation failed",
			},
			wantStatus:  http.StatusBadRequest,
			description: "should write validation error with bad request status",
		},
		{
			name: "validation error with details",
			err: &ValidationError{
				Message: "validation failed",
				Details: map[string]string{
					"title": "title is required",
				},
			},
			wantStatus:  http.StatusBadRequest,
			description: "should write validation error with details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/test", nil)

			WriteValidationError(w, req, tt.err)

			if w.Code != tt.wantStatus {
				t.Errorf("WriteValidationError() status = %v, want %v - %s", w.Code, tt.wantStatus, tt.description)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteValidationError() Content-Type = %v, want application/json", contentType)
			}

			if w.Body.Len() == 0 {
				t.Error("WriteValidationError() should write response body")
			}
		})
	}
}

func generateLongString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = 'a'
	}
	return string(result)
}

package middleware

import (
	"context"
	"net/http"

	"go-rest-api-example/internal/pkg/id"
)

type ctxRequestIDKey struct{}

var requestIDKey = ctxRequestIDKey{}

// GetRequestID returns a request ID stored in context if present.
func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

// RequestID generates a new request ID, attaches it to the request context,
// and writes it to the response header as X-Request-ID.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := id.New()

		ctx := context.WithValue(r.Context(), requestIDKey, rid)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", rid)

		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/NoroSaroyan/go-rest-api-example/internal/pkg/logger"
)

func Logging(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := GetRequestID(r.Context())

			// Create request-specific logger with request ID
			requestLogger := log.With(zap.String("request_id", reqID))

			ctx := logger.Inject(r.Context(), requestLogger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

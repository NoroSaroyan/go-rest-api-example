package app

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/NoroSaroyan/go-rest-api-example/internal/pkg/logger"
	"github.com/NoroSaroyan/go-rest-api-example/internal/service"
	v1 "github.com/NoroSaroyan/go-rest-api-example/internal/transport/http/v1"
	"github.com/NoroSaroyan/go-rest-api-example/internal/transport/middleware"
)

// NewRouter configures all HTTP routes and middleware.
func NewRouter(todoService service.TodoService, log logger.Logger) http.Handler {
	r := mux.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(log))

	// API v1
	v1Router := r.PathPrefix("/api/v1").Subrouter()

	todoHandler := v1.NewTodoHandler(todoService)
	todoHandler.RegisterRoutes(v1Router)

	// Simple healthcheck
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Error("failed to write health check response", zap.Error(err))
		}
	}).Methods("GET")

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}

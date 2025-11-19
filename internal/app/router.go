package app

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"go-rest-api-example/internal/pkg/logger"
	"go-rest-api-example/internal/service"
	v1 "go-rest-api-example/internal/transport/http/v1"
	"go-rest-api-example/internal/transport/middleware"
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
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}

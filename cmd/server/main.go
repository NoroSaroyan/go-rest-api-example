// Todo API provides a RESTful API for managing todo items
//
//	@title			Todo API
//	@version		1.0
//	@description	A REST API for managing todo items with clean architecture
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.email	support@example.com
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/api/v1
//
//	@schemes	http https
//
//	@produce	json
//	@consumes	json
package main

import (
	"context"
	"go-rest-api-example/internal/app"
	"go-rest-api-example/internal/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "go-rest-api-example/docs" // Import generated docs
)

func main() {
	log := logger.NewFromEnv()

	application, err := app.New()
	if err != nil {
		log.Fatal("failed to create application", zap.Error(err))
	}

	// graceful shutdown support
	go func() {
		if err := application.Run(); err != nil {
			log.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	log.Info("application started")

	// Wait for interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutdown signal received")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Error("failed to gracefully shutdown", zap.Error(err))
	}
}

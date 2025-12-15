package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NoroSaroyan/go-rest-api-example/internal/config"
	"github.com/NoroSaroyan/go-rest-api-example/internal/pkg/logger"
	"github.com/NoroSaroyan/go-rest-api-example/internal/postgres"
	"github.com/NoroSaroyan/go-rest-api-example/internal/repository"
	"github.com/NoroSaroyan/go-rest-api-example/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// App encapsulates the whole application state.
type App struct {
	cfg    *config.Config
	server *http.Server
	db     *pgxpool.Pool
	logger logger.Logger
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	log := logger.FromContext(ctx)
	log.Info("starting application")

	// Build DB config using functional options
	dbCfg, err := postgres.NewConfig(
		postgres.WithHost(cfg.DB.Host),
		postgres.WithPort(uint16(cfg.DB.Port)),
		postgres.WithUser(cfg.DB.User),
		postgres.WithPassword(cfg.DB.Password),
		postgres.WithDatabase(cfg.DB.Name),
		postgres.WithMaxOpenConnections(cfg.DB.MaxOpenConns),
		postgres.WithMaxIdleConnections(cfg.DB.MaxIdleConns),
		postgres.WithConnMaxLifetime(cfg.DB.ConnMaxLifetime),
		postgres.WithConnMaxIdleTime(cfg.DB.ConnMaxIdleTime),
		postgres.WithParams(make(map[string]string)),
	)
	if err != nil {
		log.Error("failed to build db config", zap.Error(err))
		return nil, fmt.Errorf("failed to build db config: %w", err)
	}

	// Open DB pool (postgres package owns pgxpool configuration + ping)
	dbpool, err := postgres.Open(ctx, dbCfg)
	if err != nil {
		log.Error("failed to connect to DB", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	log.Info("database connection established")

	// Initialize repository & service
	todoRepo := repository.NewTodoRepository(dbpool)
	todoService := service.NewTodoService(todoRepo)

	// Build router
	router := NewRouter(todoService, log)

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	return &App{
		cfg:    cfg,
		server: srv,
		db:     dbpool,
		logger: log,
	}, nil
}

// Run starts the HTTP server.
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("HTTP server listening", zap.String("port", a.cfg.App.Port))

	go func() {
		_ = a.Shutdown(ctx)
	}()

	// ListenAndServe returns http.ErrServerClosed on graceful shutdown; treat that as non-error.
	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server.
func (a *App) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	a.logger.Info("Shutting down server")

	// Stop accepting new requests / allow in-flight requests to finish.
	err := a.server.Shutdown(context.Background())

	// Close DB pool after server shutdown begins.
	a.db.Close()

	return err
}

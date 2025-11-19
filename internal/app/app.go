package app

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/NoroSaroyan/go-rest-api-example/internal/config"
	"github.com/NoroSaroyan/go-rest-api-example/internal/pkg/logger"
	"github.com/NoroSaroyan/go-rest-api-example/internal/repository"
	"github.com/NoroSaroyan/go-rest-api-example/internal/service"
)

// App encapsulates the whole application state.
type App struct {
	cfg    *config.Config
	server *http.Server
	db     *pgxpool.Pool
	logger logger.Logger
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	log := logger.New(cfg.Log.Level)
	log.Info("starting application")

	// Create DB connection with pool configuration
	dbconfig, err := pgxpool.ParseConfig(cfg.DatabaseURL())
	if err != nil {
		log.Error("failed to parse database config", zap.Error(err))
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool
	if cfg.DB.MaxOpenConns > math.MaxInt32 || cfg.DB.MaxOpenConns < 0 {
		log.Error("max open connections value out of range for int32", zap.Int("value", cfg.DB.MaxOpenConns))
		return nil, fmt.Errorf("max open connections value %d out of range for int32", cfg.DB.MaxOpenConns)
	}
	if cfg.DB.MaxIdleConns > math.MaxInt32 || cfg.DB.MaxIdleConns < 0 {
		log.Error("max idle connections value out of range for int32", zap.Int("value", cfg.DB.MaxIdleConns))
		return nil, fmt.Errorf("max idle connections value %d out of range for int32", cfg.DB.MaxIdleConns)
	}

	dbconfig.MaxConns = int32(cfg.DB.MaxOpenConns) //#nosec G115 -- bounds checked above
	dbconfig.MinConns = int32(cfg.DB.MaxIdleConns) //#nosec G115 -- bounds checked above
	dbconfig.MaxConnLifetime = cfg.DB.ConnMaxLifetime
	dbconfig.MaxConnIdleTime = cfg.DB.ConnMaxIdleTime

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbconfig)
	if err != nil {
		log.Error("failed to connect to DB", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test database connection
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Error("failed to ping database", zap.Error(err))
		dbpool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
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
func (a *App) Run() error {
	a.logger.Info("HTTP server listening", zap.String("port", a.cfg.App.Port))
	return a.server.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down server")
	a.db.Close()
	return a.server.Shutdown(ctx)
}

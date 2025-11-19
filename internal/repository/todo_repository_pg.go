// Package repository contains database implementations for the application's
// persistence layer. It provides concrete pgx-based repositories that satisfy
// consumer-defined interfaces from the service layer.
package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/NoroSaroyan/go-rest-api-example/internal/domain"
	"github.com/NoroSaroyan/go-rest-api-example/internal/pkg/logger"
)

type TodoRepositoryPg struct {
	db *pgxpool.Pool
}

// NewTodoRepository creates a new TODO repository.
func NewTodoRepository(db *pgxpool.Pool) *TodoRepositoryPg {
	return &TodoRepositoryPg{db: db}
}

// Create inserts a new todo and returns its generated ID.
func (r *TodoRepositoryPg) Create(ctx context.Context, title string) (int, error) {
	log := logger.FromContext(ctx)

	const query = `
		INSERT INTO todos (title)
		VALUES ($1)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(ctx, query, title).Scan(&id)
	if err != nil {
		log.Error("failed to insert todo", zap.Error(err))
		return 0, err
	}

	log.Info("todo created", zap.Int("id", id))
	return id, nil
}

// GetByID retrieves a todo by its ID.
func (r *TodoRepositoryPg) GetByID(ctx context.Context, id int) (*domain.Todo, error) {
	log := logger.FromContext(ctx)

	const query = `
		SELECT id, title, completed, created_at
		FROM todos
		WHERE id = $1
	`

	var t domain.Todo
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID,
		&t.Title,
		&t.Completed,
		&t.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Warn("todo not found", zap.Int("id", id))
		return nil, domain.ErrTodoNotFound
	}

	if err != nil {
		log.Error("failed to fetch todo", zap.Error(err))
		return nil, err
	}

	return &t, nil
}

// List retrieves all todos.
func (r *TodoRepositoryPg) List(ctx context.Context) ([]domain.Todo, error) {
	log := logger.FromContext(ctx)

	const query = `
		SELECT id, title, completed, created_at
		FROM todos
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		log.Error("failed to query todos", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	todos := make([]domain.Todo, 0)

	for rows.Next() {
		var t domain.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			log.Error("failed to scan todo row", zap.Error(err))
			return nil, err
		}
		todos = append(todos, t)
	}

	if rows.Err() != nil {
		log.Error("rows error", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	return todos, nil
}

// Delete removes a todo by ID.
func (r *TodoRepositoryPg) Delete(ctx context.Context, id int) error {
	log := logger.FromContext(ctx)

	const query = `
		DELETE FROM todos
		WHERE id = $1
	`

	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		log.Error("failed to delete todo", zap.Error(err))
		return err
	}

	if res.RowsAffected() == 0 {
		log.Warn("todo not found for delete", zap.Int("id", id))
		return domain.ErrTodoNotFound
	}

	log.Info("todo deleted", zap.Int("id", id))
	return nil
}

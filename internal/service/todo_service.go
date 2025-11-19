package service

import (
	"context"
	"errors"
	"go-rest-api-example/internal/domain"
	"go-rest-api-example/internal/pkg/logger"
	"go.uber.org/zap"
	"strings"
)

// TodoRepository is the contract the persistence layer must satisfy.
// The consumer (the service) owns the interface.
type TodoRepository interface {
	Create(ctx context.Context, title string) (int, error)
	GetByID(ctx context.Context, id int) (*domain.Todo, error)
	List(ctx context.Context) ([]domain.Todo, error)
	Delete(ctx context.Context, id int) error
}

// TodoService defines operations available on TODO entities.
type TodoService interface {
	Create(ctx context.Context, title string) (int, error)
	GetByID(ctx context.Context, id int) (*domain.Todo, error)
	List(ctx context.Context) ([]domain.Todo, error)
	Delete(ctx context.Context, id int) error
}

type todoService struct {
	repo TodoRepository
}

// NewTodoService constructs a new TodoService.
func NewTodoService(repo TodoRepository) TodoService {
	return &todoService{repo: repo}
}

// Create validates input and delegates todo creation to repository.
func (s *todoService) Create(ctx context.Context, title string) (int, error) {
	log := logger.FromContext(ctx)

	title = strings.TrimSpace(title)
	if title == "" {
		if log != nil {
			log.Warn("invalid empty title")
		}
		return 0, domain.ErrInvalidTitle
	}

	id, err := s.repo.Create(ctx, title)
	if err != nil {
		if log != nil {
			log.Error("failed to create todo", zap.Error(err))
		}
		return 0, err
	}

	if log != nil {
		log.Info("todo created successfully", zap.Int("id", id))
	}
	return id, nil
}

// GetByID retrieves a todo by id.
func (s *todoService) GetByID(ctx context.Context, id int) (*domain.Todo, error) {
	log := logger.FromContext(ctx)

	if id <= 0 {
		if log != nil {
			log.Warn("invalid ID provided", zap.Int("id", id))
		}
		return nil, domain.ErrTodoNotFound
	}

	t, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, domain.ErrTodoNotFound) {
		if log != nil {
			log.Warn("todo not found", zap.Int("id", id))
		}
		return nil, err
	}
	if err != nil {
		if log != nil {
			log.Error("failed to get todo", zap.Error(err))
		}
		return nil, err
	}

	return t, nil
}

// List retrieves all todos.
func (s *todoService) List(ctx context.Context) ([]domain.Todo, error) {
	log := logger.FromContext(ctx)

	todos, err := s.repo.List(ctx)
	if err != nil {
		if log != nil {
			log.Error("failed to list todos", zap.Error(err))
		}
		return nil, err
	}

	if log != nil {
		log.Info("todos fetched", zap.Int("count", len(todos)))
	}
	return todos, nil
}

// Delete removes a todo by id.
func (s *todoService) Delete(ctx context.Context, id int) error {
	log := logger.FromContext(ctx)

	if id <= 0 {
		if log != nil {
			log.Warn("invalid ID for deletion", zap.Int("id", id))
		}
		return domain.ErrTodoNotFound
	}

	err := s.repo.Delete(ctx, id)
	if errors.Is(err, domain.ErrTodoNotFound) {
		if log != nil {
			log.Warn("todo not found for deletion", zap.Int("id", id))
		}
		return err
	}
	if err != nil {
		if log != nil {
			log.Error("failed to delete todo", zap.Error(err))
		}
		return err
	}

	if log != nil {
		log.Info("todo deleted successfully", zap.Int("id", id))
	}
	return nil
}

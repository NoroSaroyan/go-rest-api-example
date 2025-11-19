package service

import (
	"context"
	"errors"
	"testing"

	"go-rest-api-example/internal/domain"
)

// MockTodoRepository implements TodoRepository for testing
type MockTodoRepository struct {
	todos  map[int]*domain.Todo
	nextID int
}

func NewMockTodoRepository() *MockTodoRepository {
	return &MockTodoRepository{
		todos:  make(map[int]*domain.Todo),
		nextID: 1,
	}
}

func (m *MockTodoRepository) Create(ctx context.Context, title string) (int, error) {
	id := m.nextID
	m.nextID++

	m.todos[id] = &domain.Todo{
		ID:    id,
		Title: title,
	}

	return id, nil
}

func (m *MockTodoRepository) GetByID(ctx context.Context, id int) (*domain.Todo, error) {
	todo, exists := m.todos[id]
	if !exists {
		return nil, domain.ErrTodoNotFound
	}
	return todo, nil
}

func (m *MockTodoRepository) List(ctx context.Context) ([]domain.Todo, error) {
	todos := make([]domain.Todo, 0, len(m.todos))
	for _, todo := range m.todos {
		todos = append(todos, *todo)
	}
	return todos, nil
}

func (m *MockTodoRepository) Delete(ctx context.Context, id int) error {
	if _, exists := m.todos[id]; !exists {
		return domain.ErrTodoNotFound
	}
	delete(m.todos, id)
	return nil
}

func TestTodoService_Create(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{
			name:  "valid title",
			title: "Test Todo",
		},
		{
			name:    "empty title",
			title:   "",
			wantErr: domain.ErrInvalidTitle,
		},
		{
			name:    "whitespace title",
			title:   "   ",
			wantErr: domain.ErrInvalidTitle,
		},
		{
			name:  "title with whitespace",
			title: "  Test Todo  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockTodoRepository()
			service := NewTodoService(repo)

			ctx := context.Background()
			id, err := service.Create(ctx, tt.title)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Create() unexpected error = %v", err)
				return
			}

			if id <= 0 {
				t.Errorf("Create() returned invalid id = %v", id)
			}

			// Verify todo was created
			todo, err := service.GetByID(ctx, id)
			if err != nil {
				t.Errorf("GetByID() failed to retrieve created todo: %v", err)
				return
			}

			expectedTitle := tt.title
			if tt.title != "" {
				expectedTitle = "Test Todo" // trimmed version
			}

			if todo.Title != expectedTitle && tt.title != "" {
				t.Errorf("Created todo title = %v, want %v", todo.Title, expectedTitle)
			}
		})
	}
}

func TestTodoService_GetByID(t *testing.T) {
	repo := NewMockTodoRepository()
	service := NewTodoService(repo)
	ctx := context.Background()

	// Create a todo first
	id, err := service.Create(ctx, "Test Todo")
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{
			name: "existing todo",
			id:   id,
		},
		{
			name:    "non-existent todo",
			id:      999,
			wantErr: domain.ErrTodoNotFound,
		},
		{
			name:    "invalid id - zero",
			id:      0,
			wantErr: domain.ErrTodoNotFound,
		},
		{
			name:    "invalid id - negative",
			id:      -1,
			wantErr: domain.ErrTodoNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := service.GetByID(ctx, tt.id)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetByID() unexpected error = %v", err)
				return
			}

			if todo == nil {
				t.Error("GetByID() returned nil todo")
				return
			}

			if todo.ID != tt.id {
				t.Errorf("GetByID() returned todo with id = %v, want %v", todo.ID, tt.id)
			}
		})
	}
}

func TestTodoService_List(t *testing.T) {
	repo := NewMockTodoRepository()
	service := NewTodoService(repo)
	ctx := context.Background()

	// Test empty list
	todos, err := service.List(ctx)
	if err != nil {
		t.Errorf("List() failed on empty repository: %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("List() returned %d todos, want 0", len(todos))
	}

	// Create some todos
	titles := []string{"Todo 1", "Todo 2", "Todo 3"}
	for _, title := range titles {
		_, err := service.Create(ctx, title)
		if err != nil {
			t.Fatalf("Failed to create test todo: %v", err)
		}
	}

	// Test list with todos
	todos, err = service.List(ctx)
	if err != nil {
		t.Errorf("List() failed: %v", err)
	}

	if len(todos) != len(titles) {
		t.Errorf("List() returned %d todos, want %d", len(todos), len(titles))
	}
}

func TestTodoService_Delete(t *testing.T) {
	repo := NewMockTodoRepository()
	service := NewTodoService(repo)
	ctx := context.Background()

	// Create a todo first
	id, err := service.Create(ctx, "Test Todo")
	if err != nil {
		t.Fatalf("Failed to create test todo: %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{
			name: "existing todo",
			id:   id,
		},
		{
			name:    "non-existent todo",
			id:      999,
			wantErr: domain.ErrTodoNotFound,
		},
		{
			name:    "invalid id - zero",
			id:      0,
			wantErr: domain.ErrTodoNotFound,
		},
		{
			name:    "invalid id - negative",
			id:      -1,
			wantErr: domain.ErrTodoNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(ctx, tt.id)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Delete() unexpected error = %v", err)
				return
			}

			// Verify todo was deleted
			_, err = service.GetByID(ctx, tt.id)
			if !errors.Is(err, domain.ErrTodoNotFound) {
				t.Errorf("GetByID() after delete should return ErrTodoNotFound, got %v", err)
			}
		})
	}
}

package v1

// CreateTodoRequest is the payload for creating a new todo.
type CreateTodoRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255" example:"Buy groceries"`
}

// TodoResponse is the JSON representation returned to clients.
type TodoResponse struct {
	ID        int    `json:"id" example:"1"`
	Title     string `json:"title" example:"Buy groceries"`
	Completed bool   `json:"completed" example:"false"`
	CreatedAt string `json:"created_at" example:"2023-01-01T12:00:00Z"`
}

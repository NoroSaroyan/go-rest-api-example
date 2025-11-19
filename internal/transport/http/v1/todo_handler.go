package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/NoroSaroyan/go-rest-api-example/internal/pkg/logger"
	"github.com/NoroSaroyan/go-rest-api-example/internal/service"
)

// TodoHandler provides HTTP endpoints for managing todos.
type TodoHandler struct {
	service service.TodoService
}

// NewTodoHandler initializes the handler.
func NewTodoHandler(s service.TodoService) *TodoHandler {
	return &TodoHandler{service: s}
}

// RegisterRoutes attaches routes to a router.
func (h *TodoHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/todos", h.create).Methods("POST")
	r.HandleFunc("/todos/{id}", h.getByID).Methods("GET")
	r.HandleFunc("/todos", h.list).Methods("GET")
	r.HandleFunc("/todos/{id}", h.delete).Methods("DELETE")
}

// CreateTodo godoc
//
//	@Summary		Create a new todo item
//	@Description	Creates a new todo item with the provided title
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			todo	body		CreateTodoRequest	true	"Todo creation request"
//	@Success		201		{object}	map[string]int		"Successfully created todo"
//	@Failure		400		{object}	ValidationError		"Validation error"
//	@Failure		500		{object}	ErrorResponse		"Internal server error"
//	@Router			/todos [post]
func (h *TodoHandler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := DecodeAndValidateJSON(r, &req); err != nil {
		if validationErr, ok := err.(*ValidationError); ok {
			WriteValidationError(w, r, validationErr)
			return
		}
		WriteError(w, r, NewValidationError("invalid request body"))
		return
	}

	id, err := h.service.Create(r.Context(), req.Title)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSONSafe(w, r, http.StatusCreated, map[string]any{"id": id})
}

// GetTodoByID godoc
//
//	@Summary		Get a todo item by ID
//	@Description	Retrieves a specific todo item by its ID
//	@Tags			todos
//	@Produce		json
//	@Param			id	path		int	true	"Todo ID"
//	@Success		200	{object}	TodoResponse		"Successfully retrieved todo"
//	@Failure		400	{object}	ErrorResponse		"Invalid ID parameter"
//	@Failure		404	{object}	ErrorResponse		"Todo not found"
//	@Failure		500	{object}	ErrorResponse		"Internal server error"
//	@Router			/todos/{id} [get]
func (h *TodoHandler) getByID(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if log != nil {
			log.Warn("invalid ID parameter", zap.String("param", idStr))
		}
		WriteError(w, r, NewValidationError("invalid id parameter"))
		return
	}

	t, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	resp := TodoResponse{
		ID:        t.ID,
		Title:     t.Title,
		Completed: t.Completed,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
	}

	WriteJSONSafe(w, r, http.StatusOK, resp)
}

// ListTodos godoc
//
//	@Summary		List all todo items
//	@Description	Retrieves a list of all todo items
//	@Tags			todos
//	@Produce		json
//	@Success		200	{array}		TodoResponse	"Successfully retrieved todos"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/todos [get]
func (h *TodoHandler) list(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.List(r.Context())
	if err != nil {
		WriteError(w, r, err)
		return
	}

	resp := make([]TodoResponse, 0, len(todos))
	for _, t := range todos {
		resp = append(resp, TodoResponse{
			ID:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		})
	}

	WriteJSONSafe(w, r, http.StatusOK, resp)
}

// DeleteTodo godoc
//
//	@Summary		Delete a todo item
//	@Description	Deletes a specific todo item by its ID
//	@Tags			todos
//	@Param			id	path	int	true	"Todo ID"
//	@Success		204	"Successfully deleted todo"
//	@Failure		400	{object}	ErrorResponse	"Invalid ID parameter"
//	@Failure		404	{object}	ErrorResponse	"Todo not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/todos/{id} [delete]
func (h *TodoHandler) delete(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if log != nil {
			log.Warn("invalid id", zap.String("param", idStr))
		}
		WriteError(w, r, NewValidationError("invalid id parameter"))
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

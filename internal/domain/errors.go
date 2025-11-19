package domain

import "errors"

// Domain-level errors returned by repositories and services,
// enabling transport layer to map them to proper HTTP responses.
var (
	ErrTodoNotFound = errors.New("todo not found")
	ErrInvalidTitle = errors.New("title cannot be empty")
)

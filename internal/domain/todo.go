// Package domain contains the core business entities used throughout the
// application. These models define the shape of data independent of any
// transport, database, or external representation.
package domain

import "time"

// Todo represents a single task item in the application.
type Todo struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Completed bool      `db:"completed"`
	CreatedAt time.Time `db:"created_at"`
}

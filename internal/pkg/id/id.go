// Package id provides small utilities for generating unique identifiers
// such as request IDs. These IDs are used to correlate log entries and trace
// individual HTTP requests through the system. The package is intentionally
// lightweight, exposing only a minimal API with no external dependencies.
//
// Typical usage:
//
//	id := id.New()    // generates a random 16-character hex string
//
// The generator uses crypto/rand for secure random bytes and is suitable
// for request correlation, logging, and lightweight tracing.
package id

import (
	"crypto/rand"
	"encoding/hex"
)

// New generates a random 16-character request ID.
// It uses crypto/rand for secure random bytes.
func New() string {
	b := make([]byte, 8) // 8 bytes = 16 hex chars
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

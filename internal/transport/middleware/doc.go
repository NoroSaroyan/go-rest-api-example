// Package middleware contains HTTP middleware used to augment incoming requests,
// including request ID generation and context-aware logging. These middlewares
// enrich requests with metadata and ensure structured, correlated logging
// throughout the application.
//
// Middlewares in this package are transport-specific and should not contain
// business logic or interact with repositories or services.
package middleware

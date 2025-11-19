// Package logger provides application-wide logging utilities built on top of
// Uber's zap library. It exposes a lazily initialized singleton logger
// configured according to environment variables, along with helpers for
// injecting and retrieving loggers from context.
//
// This package enables:
//   - Structured, leveled logging using zap
//   - Context-aware logging for per-request log enrichment
//   - A single global logger instance shared across the application
//
// Typical usage:
//
//	log := logger.New()
//	log.Info("starting application")
//
//	ctx := logger.Inject(context.Background(), log)
//	logFromCtx := logger.FromContext(ctx)
//
// The package is transport-agnostic and may be used in HTTP handlers,
// services, repository implementations, or background workers.
package logger

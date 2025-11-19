package logger

import (
	"context"
)

type ctxLoggerKey struct{}

var loggerKey = ctxLoggerKey{}

// Inject adds a logger instance into the context and returns the updated context.
func Inject(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

// FromContext retrieves the logger stored inside the context.
// If no logger is found, it returns nil - callers should handle this appropriately.
func FromContext(ctx context.Context) Logger {
	if log, ok := ctx.Value(loggerKey).(Logger); ok {
		return log
	}
	return nil
}

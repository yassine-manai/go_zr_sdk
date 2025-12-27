package logger

import "context"

// NoOpLogger is a logger that does nothing (useful for testing)
type NoOpLogger struct{}

// NewNoOpLogger creates a new no-op logger
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

// Debug does nothing
func (l *NoOpLogger) Debug(msg string, fields ...Field) {}

// Info does nothing
func (l *NoOpLogger) Info(msg string, fields ...Field) {}

// Warn does nothing
func (l *NoOpLogger) Warn(msg string, fields ...Field) {}

// Error does nothing
func (l *NoOpLogger) Error(msg string, fields ...Field) {}

// Trace does nothing
func (l *NoOpLogger) Trace(msg string, fields ...Field) {}

// With returns the same no-op logger
func (l *NoOpLogger) With(fields ...Field) Logger {
	return l
}

// WithContext returns the same no-op logger
func (l *NoOpLogger) WithContext(ctx context.Context) Logger {
	return l
}

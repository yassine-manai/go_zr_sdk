package logger

import (
	"context"
)

// Level represents logging level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelTrace
	LevelWarn
	LevelError
)

// String returns string representation of log level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelTrace:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

// Logger is the interface for logging
type Logger interface {
	// Debug logs a debug message
	Debug(msg string, fields ...Field)

	// Info logs an info message
	Info(msg string, fields ...Field)

	// Warn logs a warning message
	Warn(msg string, fields ...Field)

	// Error logs an error message
	Error(msg string, fields ...Field)

	// Error logs an Trace message
	Trace(msg string, fields ...Field)

	// With creates a child logger with additional fields
	With(fields ...Field) Logger

	// WithContext creates a logger with context
	WithContext(ctx context.Context) Logger
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Common field constructors
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

func Duration(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

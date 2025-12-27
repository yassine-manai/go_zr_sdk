package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// DefaultLogger is a simple logger implementation
type DefaultLogger struct {
	level       Level
	output      io.Writer
	mu          sync.Mutex
	fields      []Field
	prettyPrint bool
}

// DefaultLoggerOptions configures the default logger
type DefaultLoggerOptions struct {
	Level       Level
	Output      io.Writer
	PrettyPrint bool
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger(opts DefaultLoggerOptions) *DefaultLogger {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	return &DefaultLogger{
		level:       opts.Level,
		output:      opts.Output,
		fields:      make([]Field, 0),
		prettyPrint: opts.PrettyPrint,
	}
}

// Debug logs a debug message
func (l *DefaultLogger) Debug(msg string, fields ...Field) {
	if l.level <= LevelDebug {
		l.Log(LevelDebug, msg, fields...)
	}
}

// Info logs an info message
func (l *DefaultLogger) Info(msg string, fields ...Field) {
	if l.level <= LevelInfo {
		l.Log(LevelInfo, msg, fields...)
	}
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(msg string, fields ...Field) {
	if l.level <= LevelWarn {
		l.Log(LevelWarn, msg, fields...)
	}
}

// Error logs an error message
func (l *DefaultLogger) Error(msg string, fields ...Field) {
	if l.level <= LevelError {
		l.Log(LevelError, msg, fields...)
	}
}

func (l *DefaultLogger) Trace(msg string, fields ...Field) {
	if l.level <= LevelTrace {
		l.Log(LevelTrace, msg, fields...)
	}
}

// With creates a child logger with additional fields
func (l *DefaultLogger) With(fields ...Field) Logger {
	childFields := make([]Field, len(l.fields)+len(fields))
	copy(childFields, l.fields)
	copy(childFields[len(l.fields):], fields)

	return &DefaultLogger{
		level:       l.level,
		output:      l.output,
		fields:      childFields,
		prettyPrint: l.prettyPrint,
	}
}

// WithContext creates a logger with context (extracts request ID if available)
func (l *DefaultLogger) WithContext(ctx context.Context) Logger {
	// Extract request ID from context if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		return l.With(String("request_id", fmt.Sprintf("%v", requestID)))
	}
	return l
}

// log is the internal logging method
func (l *DefaultLogger) Log(level Level, msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Combine logger fields with message fields
	allFields := make([]Field, 0, len(l.fields)+len(fields))
	allFields = append(allFields, l.fields...)
	allFields = append(allFields, fields...)

	entry := logEntry{
		Timestamp: time.Now(),
		Level:     level.String(),
		Message:   msg,
		Fields:    fieldsToMap(allFields),
	}

	var output []byte
	var err error

	if l.prettyPrint {
		output, err = json.MarshalIndent(entry, "", "  ")
	} else {
		output, err = json.Marshal(entry)
	}

	if err != nil {
		fmt.Fprintf(l.output, "error marshaling log entry: %v\n", err)
		return
	}

	fmt.Fprintln(l.output, string(output))
}

// logEntry represents a log entry
type logEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// fieldsToMap converts fields to a map
func fieldsToMap(fields []Field) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	m := make(map[string]interface{}, len(fields))
	for _, f := range fields {
		m[f.Key] = f.Value
	}
	return m
}

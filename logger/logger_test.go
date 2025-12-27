package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	buf := &bytes.Buffer{}

	logger := NewDefaultLogger(DefaultLoggerOptions{
		Level:       LevelDebug,
		Output:      buf,
		PrettyPrint: false,
	})

	logger.Info("test message", String("key", "value"))

	output := buf.String()
	if output == "" {
		t.Error("expected log output, got empty string")
	}

	// Parse JSON output
	var entry logEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if entry.Message != "test message" {
		t.Errorf("expected message 'test message', got %s", entry.Message)
	}

	if entry.Level != "INFO" {
		t.Errorf("expected level INFO, got %s", entry.Level)
	}

	if entry.Fields["key"] != "value" {
		t.Errorf("expected field key=value, got %v", entry.Fields["key"])
	}
}

func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name        string
		loggerLevel Level
		logLevel    Level
		shouldLog   bool
	}{
		{"debug logger logs debug", LevelDebug, LevelDebug, true},
		{"info logger skips debug", LevelInfo, LevelDebug, false},
		{"info logger logs info", LevelInfo, LevelInfo, true},
		{"warn logger skips info", LevelWarn, LevelInfo, false},
		{"error logger only logs errors", LevelError, LevelWarn, false},
		{"trace logger skips info", LevelTrace, LevelInfo, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := NewDefaultLogger(DefaultLoggerOptions{
				Level:  tt.loggerLevel,
				Output: buf,
			})

			// Log at specific level
			switch tt.logLevel {
			case LevelDebug:
				logger.Debug("test")
			case LevelInfo:
				logger.Info("test")
			case LevelWarn:
				logger.Warn("test")
			case LevelError:
				logger.Error("test")
			}

			hasOutput := buf.Len() > 0
			if hasOutput != tt.shouldLog {
				t.Errorf("expected shouldLog=%v, got output=%v", tt.shouldLog, hasOutput)
			}
		})
	}
}

func TestLoggerWith(t *testing.T) {
	buf := &bytes.Buffer{}

	parentLogger := NewDefaultLogger(DefaultLoggerOptions{
		Level:  LevelInfo,
		Output: buf,
	})

	childLogger := parentLogger.With(String("service", "customer_media"))
	childLogger.Info("test message", String("action", "create"))

	output := buf.String()

	// Parse JSON
	var entry logEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if entry.Fields["service"] != "customer_media" {
		t.Errorf("expected service field from parent logger")
	}

	if entry.Fields["action"] != "create" {
		t.Errorf("expected action field from child logger")
	}
}

func TestLoggerWithContext(t *testing.T) {
	buf := &bytes.Buffer{}

	logger := NewDefaultLogger(DefaultLoggerOptions{
		Level:  LevelInfo,
		Output: buf,
	})

	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	contextLogger := logger.WithContext(ctx)

	contextLogger.Info("test message")

	output := buf.String()

	if !strings.Contains(output, "req-123") {
		t.Error("expected request_id in log output")
	}
}

func TestNoOpLogger(t *testing.T) {
	logger := NewNoOpLogger()

	// Should not panic
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")

	childLogger := logger.With(String("key", "value"))
	childLogger.Info("test")

	contextLogger := logger.WithContext(context.Background())
	contextLogger.Info("test")
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
		wantErr  bool
	}{
		{"debug", LevelDebug, false},
		{"DEBUG", LevelDebug, false},
		{"info", LevelInfo, false},
		{"warn", LevelWarn, false},
		{"warning", LevelWarn, false},
		{"error", LevelError, false},
		{"invalid", LevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			level, err := ParseLevel(tt.input)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.wantErr && level != tt.expected {
				t.Errorf("expected level %v, got %v", tt.expected, level)
			}
		})
	}
}

func TestFieldConstructors(t *testing.T) {
	fields := []Field{
		String("string", "value"),
		Int("int", 42),
		Int64("int64", int64(123)),
		Float64("float", 3.14),
		Bool("bool", true),
		Any("any", map[string]string{"key": "value"}),
	}

	if len(fields) != 6 {
		t.Errorf("expected 6 fields, got %d", len(fields))
	}

	// Test Error field
	errField := Error(nil)
	if errField.Key != "error" {
		t.Errorf("expected error field key, got %s", errField.Key)
	}
}

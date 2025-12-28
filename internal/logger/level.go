package logger

import (
	"fmt"
	"strings"
)

// ParseLevel parses a string into a Level
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	case "trace":
		return LevelTrace, nil
	default:
		return LevelInfo, fmt.Errorf("invalid log level: %s", s)
	}
}

// MustParseLevel parses a level or panics
func MustParseLevel(s string) Level {
	level, err := ParseLevel(s)
	if err != nil {
		panic(err)
	}
	return level
}

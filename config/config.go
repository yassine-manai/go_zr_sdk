package config

import (
	"errors"
	"fmt"
	"time"
)

// Config holds all configuration for the SDK
type Config struct {
	// UI Service configuration
	UI UIConfig

	// Database configuration
	DB DBConfig

	// Common settings
	Timeout time.Duration
	Logger  LoggerConfig
}

// UIConfig contains UI service settings
type UIConfig struct {
	Host               string // "https://20.0.0.50:8443"
	BasePath           string // "/CustomerMediaWebService"
	Username           string // For Basic Auth
	Password           string // For Basic Auth
	Timeout            time.Duration
	InsecureSkipVerify bool
}

// DBConfig contains database settings
type DBConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	SSLMode  bool // disable, require, verify-ca, verify-full
}

// LoggerConfig defines logging settings
type LoggerConfig struct {
	Level       string // debug, info, warn, error
	Enabled     bool
	PrettyPrint bool // For development
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {

	if err := c.UI.Validate(); err != nil {
		return fmt.Errorf("UI config validation failed: %w", err)
	}

	/*if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("DB config validation failed: %w", err)
	}
	*/

	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}

	return nil
}

// Validate checks UI configuration
func (u *UIConfig) Validate() error {
	if u.Host == "" {
		return errors.New("host is required")
	}

	if u.Username == "" {
		return errors.New("auth username is required")
	}

	if u.Password == "" {
		return errors.New("auth password is required")
	}

	return nil
}

// Validate checks DB configuration
func (d *DBConfig) Validate() error {
	if d.Host == "" {
		return errors.New("database host is required")
	}

	if d.Port <= 0 || d.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if d.Database == "" {
		return errors.New("database name is required")
	}

	if d.Username == "" {
		return errors.New("database username is required")
	}

	// Password can be empty for some auth methods

	return nil
}

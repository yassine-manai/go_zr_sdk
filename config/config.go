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
	Timeout     time.Duration
	RetryConfig RetryConfig
	Logger      LoggerConfig
}

// UIConfig contains UI service settings
type UIConfig struct {
	Host               string // "https://20.0.0.55:8443"
	BasePath           string // "/CustomerMediaWebService"
	Username           string // For Basic Auth
	Password           string // For Basic Auth
	Timeout            time.Duration
	InsecureSkipVerify bool
}

// DBConfig contains database settings
type DBConfig struct {
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	SSLMode         string // disable, require, verify-ca, verify-full
}

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxRetries     int           // Maximum number of retry attempts
	InitialBackoff time.Duration // Initial backoff duration
	MaxBackoff     time.Duration // Maximum backoff duration
	Multiplier     float64       // Backoff multiplier
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

	if err := c.RetryConfig.Validate(); err != nil {
		return fmt.Errorf("retry config validation failed: %w", err)
	}

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

	if d.MaxOpenConns < 0 {
		return errors.New("max open connections cannot be negative")
	}

	if d.MaxIdleConns < 0 {
		return errors.New("max idle connections cannot be negative")
	}

	return nil
}

// Validate checks retry configuration
func (r *RetryConfig) Validate() error {
	if r.MaxRetries < 0 {
		return errors.New("max retries cannot be negative")
	}

	if r.InitialBackoff <= 0 {
		return errors.New("initial backoff must be greater than 0")
	}

	if r.MaxBackoff < r.InitialBackoff {
		return errors.New("max backoff must be greater than or equal to initial backoff")
	}

	if r.Multiplier <= 0 {
		return errors.New("multiplier must be greater than 0")
	}

	return nil
}

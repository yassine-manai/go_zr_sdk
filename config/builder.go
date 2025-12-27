package config

import "time"

// Builder provides a fluent interface for building Config
type Builder struct {
	config *Config
}

// NewBuilder creates a new config builder with defaults
func NewBuilder() *Builder {
	return &Builder{
		config: DefaultConfig(),
	}
}

// WithUIConfig sets UI configuration
func (b *Builder) WithUIConfig(host, username, password string) *Builder {
	b.config.UI.Host = host
	b.config.UI.Username = username
	b.config.UI.Password = password
	return b
}

// WithUITimeout sets UI-specific timeout
func (b *Builder) WithUITimeout(timeout time.Duration) *Builder {
	b.config.UI.Timeout = timeout
	return b
}

// WithDBConfig sets database configuration
func (b *Builder) WithDBConfig(host, database, username, password string, port int) *Builder {
	b.config.DB.Host = host
	b.config.DB.Port = port
	b.config.DB.Database = database
	b.config.DB.Username = username
	b.config.DB.Password = password
	return b
}

// WithDBConnectionPool sets database connection pool settings
func (b *Builder) WithDBConnectionPool(maxOpen, maxIdle int, maxLifetime time.Duration) *Builder {
	b.config.DB.MaxOpenConns = maxOpen
	b.config.DB.MaxIdleConns = maxIdle
	b.config.DB.ConnMaxLifetime = maxLifetime
	return b
}

// WithTimeout sets global timeout
func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	b.config.Timeout = timeout
	return b
}

// WithRetry sets retry configuration
func (b *Builder) WithRetry(maxRetries int, initialBackoff, maxBackoff time.Duration) *Builder {
	b.config.RetryConfig.MaxRetries = maxRetries
	b.config.RetryConfig.InitialBackoff = initialBackoff
	b.config.RetryConfig.MaxBackoff = maxBackoff
	return b
}

// WithLogger sets logger configuration
func (b *Builder) WithLogger(level string, enabled bool) *Builder {
	b.config.Logger.Level = level
	b.config.Logger.Enabled = enabled
	return b
}

// Build validates and returns the configuration
func (b *Builder) Build() (*Config, error) {
	if err := b.config.Validate(); err != nil {
		return nil, err
	}
	return b.config, nil
}

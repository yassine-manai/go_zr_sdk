package config

import "time"

// Builder provides a fluent interface for building Config
type Builder struct {
	config *Config
}

// NewBuilder creates a new config builder with defaults
func NewBuilder() *Builder {
	return &Builder{
		config: &Config{},
	}
}

// WithUIConfig sets UI configuration
func (b *Builder) WithUIConfig(host, username, password string, insecureTLS bool) *Builder {
	b.config.UI.Host = host
	b.config.UI.Username = username
	b.config.UI.Password = password
	b.config.UI.InsecureSkipVerify = insecureTLS
	return b
}

// WithUITimeout sets UI-specific timeout
func (b *Builder) WithUITimeout(timeout time.Duration) *Builder {
	b.config.UI.Timeout = timeout
	return b
}

// WithDBConfig sets database configuration
func (b *Builder) WithDBConfig(host string, port int, database, username, password string) *Builder {
	b.config.DB.Host = host
	b.config.DB.Port = port
	b.config.DB.Database = database
	b.config.DB.Username = username
	b.config.DB.Password = password
	return b
}

// WithTimeout sets global timeout
func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	b.config.Timeout = timeout
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

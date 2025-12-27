package config

import "time"

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Timeout: 30 * time.Second,

		RetryConfig: RetryConfig{
			MaxRetries:     3,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     30 * time.Second,
			Multiplier:     2.0,
		},

		Logger: LoggerConfig{
			Level:       "info",
			Enabled:     true,
			PrettyPrint: false,
		},

		DB: DBConfig{
			Port:            5432,
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
			SSLMode:         "require",
		},
	}
}

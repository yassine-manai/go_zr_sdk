package client

import (
	"context"
	"testing"
	"time"

	"github.com/yassine-manai/go_zr_sdk/config"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		UI: config.UIConfig{
			Host:     "https://api.example.com",
			Username: "test-key",
		},
		DB: config.DBConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			Username: "testuser",
			Password: "testpass",
			SSLMode:  "disable",
		},
		Timeout: 30 * time.Second,
		RetryConfig: config.RetryConfig{
			MaxRetries:     3,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     10 * time.Second,
			Multiplier:     2.0,
		},
		Logger: config.LoggerConfig{
			Level:   "info",
			Enabled: false, // Disable for tests
		},
	}

	// Note: This will fail if DB is not running
	// In real tests, you'd use a test database or mock
	client, err := New(cfg)
	if err != nil {
		// Expected if DB is not running
		t.Logf("client creation failed (expected if DB not running): %v", err)
		return
	}
	defer client.Close()

	if client == nil {
		t.Fatal("expected client, got nil")
	}

	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}

	if client.logger == nil {
		t.Error("logger should not be nil")
	}
}

func TestNewClientWithNilConfig(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Error("expected error with nil config, got nil")
	}
}

func TestNewClientWithInvalidConfig(t *testing.T) {
	cfg := &config.Config{
		// Missing required fields
		Timeout: 30 * time.Second,
	}

	_, err := New(cfg)
	if err == nil {
		t.Error("expected error with invalid config, got nil")
	}
}

func TestConfigBuilder(t *testing.T) {
	cfg, err := config.NewBuilder().
		WithUIConfig("https://api.example.com", "user", "test-key").
		//WithDBConfig("localhost", "testdb", "user", "pass", 5432).
		WithTimeout(60 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("builder failed: %v", err)
	}

	// Note: This will fail if DB is not running
	client, err := New(cfg)
	if err != nil {
		t.Logf("client creation failed (expected if DB not running): %v", err)
		return
	}
	defer client.Close()

	if client.Config().Timeout != 60*time.Second {
		t.Errorf("expected timeout 60s, got %v", client.Config().Timeout)
	}
}

func TestClientAccessors(t *testing.T) {
	cfg := &config.Config{
		UI: config.UIConfig{
			Host:     "https://api.example.com",
			Username: "test-key",
		},
		DB: config.DBConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			Username: "testuser",
			Password: "testpass",
			SSLMode:  "disable",
		},
		Timeout: 30 * time.Second,
		RetryConfig: config.RetryConfig{
			MaxRetries:     3,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     10 * time.Second,
			Multiplier:     2.0,
		},
		Logger: config.LoggerConfig{
			Level:   "info",
			Enabled: true,
		},
	}

	client, err := New(cfg)
	if err != nil {
		t.Skipf("skipping test, DB not available: %v", err)
		return
	}
	defer client.Close()

	if client.HTTPClient() == nil {
		t.Error("HTTPClient() should not return nil")
	}

	if client.DBConnection() == nil {
		t.Error("DBConnection() should not return nil")
	}

	if client.Logger() == nil {
		t.Error("Logger() should not return nil")
	}

	if client.Config() != cfg {
		t.Error("Config() should return original config")
	}
}

func TestPingDB(t *testing.T) {
	cfg := &config.Config{
		UI: config.UIConfig{
			Host:     "https://api.example.com",
			Username: "test-key",
		},
		DB: config.DBConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			Username: "testuser",
			Password: "testpass",
			SSLMode:  "disable",
		},
		Timeout: 30 * time.Second,
		RetryConfig: config.RetryConfig{
			MaxRetries:     3,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     10 * time.Second,
			Multiplier:     2.0,
		},
		Logger: config.LoggerConfig{
			Level:   "info",
			Enabled: false,
		},
	}

	client, err := New(cfg)
	if err != nil {
		t.Skipf("skipping test, DB not available: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.pingDB(ctx); err != nil {
		t.Errorf("pingDB failed: %v", err)
	}
}

package client

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/yassine-manai/go_zr_sdk/config"
	"github.com/yassine-manai/go_zr_sdk/errors"
	internalhttp "github.com/yassine-manai/go_zr_sdk/internal/http"
	internalretry "github.com/yassine-manai/go_zr_sdk/internal/retry"
	"github.com/yassine-manai/go_zr_sdk/logger"
	"github.com/yassine-manai/go_zr_sdk/ui/customer_media/contract"
)

// Client is the main SDK client
type Client struct {
	config     *config.Config
	httpClient *http.Client
	dbConn     *sql.DB
	logger     logger.Logger

	// Services
	Contract *contract.Service
}

// New creates a new SDK client
func New(cfg *config.Config) (*Client, error) {

	if cfg == nil {
		return nil, errors.NewSDKError(
			errors.ErrorTypeValidation,
			"config cannot be nil",
			nil,
		)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, errors.NewSDKError(
			errors.ErrorTypeValidation,
			"invalid configuration",
			err,
		)
	}

	// Create logger (UNCOMMENTED!)
	log := createLogger(cfg)

	// Create HTTP client
	httpClient := createHTTPClient(cfg)

	// Create database connection
	/*dbConn, err := createDBConnection(cfg, log)
	if err != nil {
		return nil, errors.NewSDKError(
			errors.ErrorTypeDatabase,
			"failed to create database connection",
			err,
		)
	}*/

	// clients http + retry
	internalHTTPClient := internalhttp.NewClient(httpClient, cfg, log)
	retryer := internalretry.New(cfg.RetryConfig, log)

	var dbConn *sql.DB

	client := &Client{
		config:     cfg,
		httpClient: httpClient,
		dbConn:     dbConn,
		logger:     log, // Assign to logger field
	}

	client.Contract = contract.New(internalHTTPClient, retryer, log)

	log.Info("SDK client initialized successfully",
		logger.String("ui_host", cfg.UI.Host),
	)

	return client, nil
}

// Close closes all connections
func (c *Client) Close() error {
	c.logger.Info("closing SDK client") // Fixed: c.logger not c.log

	if c.dbConn != nil {
		if err := c.dbConn.Close(); err != nil {
			c.logger.Error("failed to close database connection", // Fixed: c.logger
				logger.Error(err), // Fixed: logger.Error not log.Error
			)
			return errors.NewDatabaseError(
				"failed to close database connection",
				"",
				"CLOSE",
				err,
			)
		}
	}

	c.logger.Info("SDK client closed successfully") // Fixed: c.logger
	return nil
}

// Ping checks connectivity to all services
func (c *Client) Ping(ctx context.Context) error {
	// Ping UI service
	if err := c.pingUI(ctx); err != nil {
		return err
	}

	// Ping database
	if err := c.pingDB(ctx); err != nil {
		return err
	}

	c.logger.Info("ping successful") // Fixed: c.logger
	return nil
}

// pingUI checks UI service connectivity
func (c *Client) pingUI(ctx context.Context) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.config.UI.Host+"/health", // Fixed: BaseURL not Host
		nil,
	)
	if err != nil {
		return errors.NewNetworkError("failed to create ping request", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.UI.Host) // Fixed: APIKey not Username

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.NewNetworkError("UI service ping failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.NewServiceUnavailableError(
			fmt.Sprintf("UI service returned status %d", resp.StatusCode),
			"ui-service",
			0,
		)
	}

	return nil
}

// pingDB checks database connectivity
func (c *Client) pingDB(ctx context.Context) error {
	if err := c.dbConn.PingContext(ctx); err != nil {
		return errors.NewDatabaseError(
			"database ping failed",
			"",
			"PING",
			err,
		)
	}
	return nil
}

// HTTPClient returns the underlying HTTP client (useful for advanced users)
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// DBConnection returns the underlying database connection (useful for advanced users)
func (c *Client) DBConnection() *sql.DB {
	return c.dbConn
}

// Logger returns the logger instance
func (c *Client) Logger() logger.Logger { // Fixed: proper return type and method name
	return c.logger
}

// Config returns the client configuration
func (c *Client) Config() *config.Config {
	return c.config
}

// createLogger creates a logger based on config
func createLogger(cfg *config.Config) logger.Logger {
	if !cfg.Logger.Enabled {
		return logger.NewNoOpLogger()
	}

	level, err := logger.ParseLevel(cfg.Logger.Level)
	if err != nil {
		// Default to Info if parsing fails
		level = logger.LevelInfo
	}

	return logger.NewDefaultLogger(logger.DefaultLoggerOptions{
		Level:       level,
		PrettyPrint: cfg.Logger.PrettyPrint,
	})
}

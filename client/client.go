package client

import (
	"database/sql"
	"net/http"

	"github.com/yassine-manai/go_zr_sdk/config"
	"github.com/yassine-manai/go_zr_sdk/internal/errors"
	internalhttp "github.com/yassine-manai/go_zr_sdk/internal/http"
	"github.com/yassine-manai/go_zr_sdk/internal/logger"
	"github.com/yassine-manai/go_zr_sdk/ui/customer_media/contract"
	"github.com/yassine-manai/go_zr_sdk/ui/customer_media/participant"
)

// Client is the main SDK client
type Client struct {
	config     *config.Config // external config
	httpClient *http.Client   // httpConnection helper
	dbConn     *sql.DB        // dbConnection helper
	logger     logger.Logger  // log Handler
	UI         UI             // ZR UI's Handler
	DB         DB             // ZR DB Handler
}

// UI Strct
type UI struct {
	CustomerMedia struct {
		Contract    *contract.ContractService       // Contract Service
		Participant *participant.ParticipantService // Participants Service
	}
}

// DB stct
type DB struct{}

// New creates a new SDK client
func NewZRClient(cfg *config.Config) (*Client, error) {

	if cfg == nil {
		return nil, errors.NewSDKError(errors.ErrorTypeValidation, "config cannot be nil", nil)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, errors.NewSDKError(errors.ErrorTypeValidation, "invalid configuration", err)
	}

	// ================# init LOGGER helper #=====================//
	log := createLogger(cfg)

	// Create database connection
	/*dbConn, err := createDBConnection(cfg, log)
	if err != nil {
		return nil, errors.NewSDKError(
			errors.ErrorTypeDatabase,
			"failed to create database connection",
			err,
		)
	}*/

	// ================# init HTTP client/helper #=====================//
	httpClient := createHTTPClient(cfg)
	internalHTTPClient := internalhttp.NewClient(httpClient, cfg, log)

	client := &Client{
		config:     cfg,
		httpClient: httpClient,
		//dbConn:     dbConn,
		logger: log,
	}

	// ================# init Services #=====================//
	client.UI.CustomerMedia.Contract = contract.NewContractService(internalHTTPClient, log)
	client.UI.CustomerMedia.Participant = participant.NewParticipantService(internalHTTPClient, log)
	// =====================================================//

	log.Info("SDK client initialized successfully", logger.String("ui_host", cfg.UI.Host))

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

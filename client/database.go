package client

import (
	"database/sql"
	"fmt"

	"github.com/yourcompany/thirdparty-sdk/config"
	"github.com/yourcompany/thirdparty-sdk/logger"
)

// createDBConnection creates and configures database connection
func createDBConnection(cfg *config.Config, log logger.Logger) (*sql.DB, error) {
	// Build connection string
	connStr := buildConnectionString(cfg.DB)

	log.Debug("connecting to database",
		logger.String("host", cfg.DB.Host),
		logger.Int("port", cfg.DB.Port),
		logger.String("database", cfg.DB.Database),
	)

	// Open connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Info("database connection established")

	return db, nil
}

// buildConnectionString constructs PostgreSQL connection string
func buildConnectionString(cfg config.DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)
}

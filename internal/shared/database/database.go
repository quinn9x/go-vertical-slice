// Package database provides database connection and access functionality.
package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/quinn9x/go-vertical-slice/internal/shared/config"
)

const (
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 10
	defaultConnMaxLifetime = 30 * time.Minute
)

// DB wraps a GORM database connection.
type DB struct {
	*gorm.DB
}

// New opens a database connection using the provided connection string.
func New(cfg *config.DatabaseConfig) (*DB, error) {
	db, err := gorm.Open(
		postgres.Open(cfg.DSN()),
		&gorm.Config{},
	)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql database: %w", err)
	}

	sqlDB.SetMaxOpenConns(defaultMaxOpenConns)
	sqlDB.SetMaxIdleConns(defaultMaxIdleConns)
	sqlDB.SetConnMaxLifetime(defaultConnMaxLifetime)

	return &DB{DB: db}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Ping checks whether the database connection is available.
func (db *DB) Ping() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

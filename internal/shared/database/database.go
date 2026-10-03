// Package database provides database connection and access functionality.
package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const databaseDirPermissions = 0o750

// Database wraps a GORM database connection.
type Database struct {
	*gorm.DB
}

// New opens a database connection using the provided connection string.
func New(connStr string) (*Database, error) {
	dir := filepath.Dir(connStr)

	if err := os.MkdirAll(dir, databaseDirPermissions); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return &Database{DB: db}, nil
}

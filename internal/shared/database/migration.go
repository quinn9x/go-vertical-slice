package database

import (
	"fmt"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
)

// Migrate runs all database migrations.
func Migrate(db *DB) error {
	if err := db.AutoMigrate(&domain.Product{}); err != nil {
		return fmt.Errorf("run database migrations: %w", err)
	}

	return nil
}

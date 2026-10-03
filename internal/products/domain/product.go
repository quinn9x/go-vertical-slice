// Package domain defines the core business logic and entities for the product feature.
package domain

import "time"

// Product represents a product in the system.
type Product struct {
	ID        string    `gorm:"primaryKey;size:36"`
	Name      string    `gorm:"not null;size:100"`
	Price     float64   `gorm:"not null"`
	Version   int64     `gorm:"not null;default:1"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

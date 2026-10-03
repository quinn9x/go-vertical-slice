// Package infra provides the implementation of the product repository using GORM for database interactions.
package infra

import (
	"context"

	"gorm.io/gorm"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
)

// ProductRepository implements the Repository interface for managing products in the database.
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepository with the provided GORM database connection.
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create inserts a new product into the database.
func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

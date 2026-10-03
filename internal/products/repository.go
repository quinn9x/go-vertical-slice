// Package products defines the product-related repositories and data access logic.
package products

import (
	"context"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
)

// Repository defines the interface for the product repository.
type Repository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	List(ctx context.Context, options ListOptions) ([]domain.Product, int64, error)
	Update(ctx context.Context, product *domain.Product, expectedVersion int64) error
	Delete(ctx context.Context, id string) error
}

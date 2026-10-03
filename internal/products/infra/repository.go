// Package infra provides the implementation of the product repository using GORM for database interactions.
package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/quinn9x/go-vertical-slice/internal/products"
	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
)

const sortByName = "name"

const sortByPrice = "price"

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

// GetByID retrieves a product from the database by its ID.
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var product domain.Product

	err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.NewNotFound("product not found")
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

// List retrieves a paginated list of products from the database based on the provided pagination parameters.
func (r *ProductRepository) List(ctx context.Context, options products.ListOptions) ([]domain.Product, int64, error) {
	var items []domain.Product

	var total int64

	query := r.db.
		WithContext(ctx).
		Model(&domain.Product{})

	if options.Search != "" {
		search := "%" + options.Search + "%"

		query = query.Where(
			"name LIKE ?",
			search,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order(
			sortColumn(options.SortBy) +
				" " +
				sortDirection(options.SortOrder),
		).
		Offset(options.Pagination.Offset()).
		Limit(options.Pagination.PageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Update modifies an existing product in the database based on its ID.
func (r *ProductRepository) Update(ctx context.Context, product *domain.Product, expectedVersion int64) error {
	result := r.db.
		WithContext(ctx).
		Model(&domain.Product{}).
		Where(
			"id = ? AND version = ?",
			product.ID,
			expectedVersion,
		).
		Updates(map[string]any{
			"name":    product.Name,
			"price":   product.Price,
			"version": gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.NewConflict(
			"Product was modified by another request",
		)
	}

	return nil
}

// Delete removes a product from the database based on its ID.
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Product{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.NewNotFound("product not found")
	}

	return nil
}

func sortColumn(sortBy string) string {
	switch sortBy {
	case sortByName:
		return sortByName

	case sortByPrice:
		return sortByPrice

	default:
		return "created_at"
	}
}

func sortDirection(sortOrder string) string {
	switch sortOrder {
	case "asc":
		return "ASC"

	default:
		return "DESC"
	}
}

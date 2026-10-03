package list

import (
	"context"

	"github.com/quinn9x/go-vertical-slice/internal/products"
	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	"github.com/quinn9x/go-vertical-slice/internal/shared/pagination"
)

// ProductResponse represents the response structure for a product in the listing.
type ProductResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Handler handles the listing of products.
type Handler struct {
	repository products.Repository
}

// NewHandler creates a new instance of Handler with the provided product repository.
func NewHandler(repository products.Repository) *Handler {
	return &Handler{repository: repository}
}

// Handle retrieves a paginated list of products based on the provided pagination parameters.
func (h *Handler) Handle(ctx context.Context, query ListerProductsQuery) (pagination.Result[ProductResponse], error) {
	items, total, err := h.repository.List(
		ctx,
		products.ListOptions{
			Pagination: query.Pagination,
			Search:     query.Search,
			SortBy:     query.SortBy,
			SortOrder:  query.SortOrder,
		})
	if err != nil {
		return pagination.Result[ProductResponse]{}, err
	}

	responses := make([]ProductResponse, 0, len(items))

	for _, product := range items {
		responses = append(
			responses,
			toResponse(&product),
		)
	}

	return pagination.Result[ProductResponse]{
		Items:    responses,
		Total:    total,
		PageSize: query.Pagination.PageSize,
		Page:     query.Pagination.Page,
	}, nil
}

// toResponse converts a domain.Product to a ProductResponse.
func toResponse(product *domain.Product) ProductResponse {
	return ProductResponse{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}
}

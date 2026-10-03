package get

import (
	"context"

	"github.com/quinn9x/go-vertical-slice/internal/products"
	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
)

// ProductResponse represents the response structure for a product.
type ProductResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// Handler handles the product retrieval logic.
type Handler struct {
	repository products.Repository
}

// NewHandler creates a new instance of Handler with the provided product repository.
func NewHandler(repository products.Repository) *Handler {
	return &Handler{repository: repository}
}

// Handle retrieves a product by its ID and returns the corresponding ProductResponse.
func (h *Handler) Handle(ctx context.Context, query GetterProductQuery) (ProductResponse, error) {
	product, err := h.repository.GetByID(ctx, query.ID)
	if err != nil {
		return ProductResponse{}, err
	}

	return toResponse(product), nil
}

func toResponse(product *domain.Product) ProductResponse {
	return ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

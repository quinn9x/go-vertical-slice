package update

import (
	"context"

	"github.com/quinn9x/go-vertical-slice/internal/products"
	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

// ProductResponse represents the response structure for a product after an update operation.
type ProductResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Version   int64   `json:"version"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// Handler is responsible for handling the update product command.
type Handler struct {
	repository products.Repository
	validator  *validation.Validator
}

// NewHandler creates a new instance of Handler with the provided product repository.
func NewHandler(validator *validation.Validator, repository products.Repository) *Handler {
	return &Handler{
		repository: repository,
		validator:  validator,
	}
}

// Handle processes the UpdateProductCommand, updating the product in the repository and returning a ProductResponse.
func (h *Handler) Handle(ctx context.Context, command UpdaterProductCommand) (ProductResponse, error) {
	if err := h.validator.Struct(command); err != nil {
		return ProductResponse{}, err
	}

	product, err := h.repository.GetByID(ctx, command.ID)
	if err != nil {
		return ProductResponse{}, err
	}

	product.Name = command.Name
	product.Price = command.Price

	if err := h.repository.Update(ctx, product, command.Version); err != nil {
		return ProductResponse{}, err
	}

	product.Version = command.Version + 1

	return toResponse(product), nil
}

func toResponse(product *domain.Product) ProductResponse {
	return ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Version:   product.Version,
		CreatedAt: product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

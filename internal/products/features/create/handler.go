package create

import (
	"context"

	"github.com/google/uuid"

	"github.com/quinn9x/go-vertical-slice/internal/products"
	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

// ProductResponse represents a product returned by the API.
type ProductResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Handler handles creator product requests.
type Handler struct {
	validator  *validation.Validator
	repository products.Repository
}

// NewHandler creates a new product creator handler.
func NewHandler(validator *validation.Validator, repository products.Repository) *Handler {
	return &Handler{
		validator:  validator,
		repository: repository,
	}
}

// Handle creates a product response from the given command.
func (h *Handler) Handle(ctx context.Context, command CreatorProductCommand) (ProductResponse, error) {
	if err := h.validator.Struct(command); err != nil {
		return ProductResponse{}, err
	}

	product := &domain.Product{
		ID:    uuid.NewString(),
		Name:  command.Name,
		Price: command.Price,
	}

	if err := h.repository.Create(ctx, product); err != nil {
		return ProductResponse{}, err
	}

	return ProductResponse{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}

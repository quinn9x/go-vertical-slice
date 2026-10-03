package create

import (
	"github.com/google/uuid"

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
	validator *validation.Validator
}

// NewHandler creates a new product creator handler.
func NewHandler(validator *validation.Validator) *Handler {
	return &Handler{
		validator: validator,
	}
}

// Handle creates a product response from the given command.
func (h *Handler) Handle(command CreatorProductCommand) (ProductResponse, error) {
	if err := h.validator.Struct(command); err != nil {
		return ProductResponse{}, err
	}

	return ProductResponse{
		ID:    uuid.NewString(),
		Name:  command.Name,
		Price: command.Price,
	}, nil
}

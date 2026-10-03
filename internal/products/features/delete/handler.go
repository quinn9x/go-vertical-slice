package deleter

import (
	"context"

	"github.com/quinn9x/go-vertical-slice/internal/products"
)

// Handler handles the deletion of products.
type Handler struct {
	repository products.Repository
}

// NewHandler creates a new instance of Handler with the provided product repository.
func NewHandler(repository products.Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

// Handle processes the DeleteProductCommand to delete a product by its ID.
func (h *Handler) Handle(
	ctx context.Context,
	command DeleteProductCommand,
) error {
	return h.repository.Delete(ctx, command.ID)
}

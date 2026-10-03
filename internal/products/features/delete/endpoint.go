package deleter

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Endpoint represents the HTTP endpoint for deleting products.
type Endpoint struct {
	handler *Handler
}

// NewEndpoint creates a new instance of Endpoint with the provided handler.
func NewEndpoint(handler *Handler) *Endpoint {
	return &Endpoint{
		handler: handler,
	}
}

// Handle processes the HTTP request to delete a product by its ID.
func (e *Endpoint) Handle(c *echo.Context) error {
	command := DeleteProductCommand{
		ID: c.Param("id"),
	}

	if err := e.handler.Handle(
		c.Request().Context(),
		command,
	); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

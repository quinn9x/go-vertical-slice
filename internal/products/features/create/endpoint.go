package create

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Endpoint is the HTTP endpoint for creating a product.
type Endpoint struct {
	handler *Handler
}

// NewEndpoint creates a new Endpoint instance with the given handler.
func NewEndpoint(handler *Handler) *Endpoint {
	return &Endpoint{
		handler: handler,
	}
}

// Handle processes the incoming HTTP request to create a product.
func (e *Endpoint) Handle(c *echo.Context) error {
	var command CreatorProductCommand

	if err := c.Bind(&command); err != nil {
		return err
	}

	result, err := e.handler.Handle(c.Request().Context(), command)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, result)
}

package update

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Endpoint is responsible for handling the update product endpoint.
type Endpoint struct {
	handler *Handler
}

// NewEndpoint creates a new instance of Endpoint with the provided handler.
func NewEndpoint(handler *Handler) *Endpoint {
	return &Endpoint{
		handler: handler,
	}
}

// Handle processes the HTTP request for updating a product, binding the request data to an UpdateProductCommand and invoking the handler.
func (e *Endpoint) Handle(c *echo.Context) error {
	var command UpdaterProductCommand
	if err := c.Bind(&command); err != nil {
		return err
	}

	command.ID = c.Param("id")

	result, err := e.handler.Handle(c.Request().Context(), command)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

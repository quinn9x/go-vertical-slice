package get

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Endpoint represents the HTTP endpoint for product retrieval.
type Endpoint struct {
	handler *Handler
}

// NewEndpoint creates a new instance of Endpoint with the provided handler.
func NewEndpoint(handler *Handler) *Endpoint {
	return &Endpoint{handler: handler}
}

// Handle processes the HTTP request to retrieve a product by its ID and returns the corresponding JSON response.
func (e *Endpoint) Handle(c *echo.Context) error {
	query := GetterProductQuery{ID: c.Param("id")}

	result, err := e.handler.Handle(c.Request().Context(), query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

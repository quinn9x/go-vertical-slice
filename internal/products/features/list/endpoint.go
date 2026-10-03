package list

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
)

// Endpoint is the HTTP endpoint for listing products.
type Endpoint struct {
	handler *Handler
}

// NewEndpoint creates a new Endpoint instance with the given handler.
func NewEndpoint(handler *Handler) *Endpoint {
	return &Endpoint{handler: handler}
}

// Handle processes the HTTP request for listing products and returns the response.
func (e *Endpoint) Handle(c *echo.Context) error {
	page, err := parseQueryInt(c.QueryParam("page"))
	if err != nil {
		return err
	}

	pageSize, err := parseQueryInt(c.QueryParam("pageSize"))
	if err != nil {
		return err
	}

	query := NewQuery(
		page,
		pageSize,
		c.QueryParam("search"),
		c.QueryParam("sortBy"),
		c.QueryParam("sortOrder"),
	)

	err = query.Pagination.Validate()
	if err != nil {
		return err
	}

	err = validateQuery(query)
	if err != nil {
		return err
	}

	result, err := e.handler.Handle(c.Request().Context(), query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func parseQueryInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, apperrors.NewBadRequest(
			"Invalid pagination parameter",
		)
	}

	return result, nil
}

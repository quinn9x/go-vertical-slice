// Package pagination provides constants and types for pagination.
package pagination

import apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"

// Constants for default pagination values.
const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Params represents the pagination parameters.
type Params struct {
	Page     int
	PageSize int
}

// New creates pagination parameters using defaults for zero values.
func New(page, pageSize int) Params {
	if page == 0 {
		page = DefaultPage
	}

	if pageSize == 0 {
		pageSize = DefaultPageSize
	}

	return Params{
		Page:     page,
		PageSize: pageSize,
	}
}

// Validate validates the pagination parameters.
func (p Params) Validate() error {
	if p.Page < 1 {
		return apperrors.NewBadRequest("Page must be greater than 0")
	}

	if p.PageSize < 1 || p.PageSize > MaxPageSize {
		return apperrors.NewBadRequest(
			"Page size must be between 1 and 100",
		)
	}

	return nil
}

// Offset returns the zero-based offset for the current page.
func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Result represents a paginated response.
type Result[T any] struct {
	Items    []T   `json:"items"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

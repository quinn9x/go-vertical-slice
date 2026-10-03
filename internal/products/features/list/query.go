// Package list provides the listing functionality for products.
package list

import (
	"strings"

	"github.com/quinn9x/go-vertical-slice/internal/shared/pagination"
)

// ListerProductsQuery represents the query parameters for listing products.
type ListerProductsQuery struct {
	Pagination pagination.Params
	Search     string
	SortBy     string
	SortOrder  string
}

// NewQuery creates a new ListerProductsQuery instance with the provided parameters.
func NewQuery(
	page int,
	pageSize int,
	search string,
	sortBy string,
	sortOrder string,
) ListerProductsQuery {
	if sortBy == "" {
		sortBy = DefaultSortBy
	}

	if sortOrder == "" {
		sortOrder = DefaultSortOrder
	}

	return ListerProductsQuery{
		Pagination: pagination.New(page, pageSize),
		Search:     strings.TrimSpace(search),
		SortBy:     sortBy,
		SortOrder:  strings.ToLower(sortOrder),
	}
}

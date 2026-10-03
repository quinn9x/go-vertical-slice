package products

import "github.com/quinn9x/go-vertical-slice/internal/shared/pagination"

// ListOptions represents the options for listing products.
type ListOptions struct {
	Pagination pagination.Params
	Search     string
	SortBy     string
	SortOrder  string
}

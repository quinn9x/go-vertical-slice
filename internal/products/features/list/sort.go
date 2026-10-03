package list

// ListProductsQuery represents the query parameters for listing products.
const (
	SortByCreatedAt = "createdAt"
	SortByName      = "name"
	SortByPrice     = "price"

	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"

	DefaultSortBy    = SortByCreatedAt
	DefaultSortOrder = SortOrderDesc
)

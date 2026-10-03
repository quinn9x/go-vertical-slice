package list_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	prodList "github.com/quinn9x/go-vertical-slice/internal/products/features/list"
)

func TestNewQuery_UsesDefaultValues(t *testing.T) {
	t.Parallel()

	query := prodList.NewQuery(
		0,
		0,
		"",
		"",
		"",
	)

	assert.Equal(t, 1, query.Pagination.Page)
	assert.Equal(t, 10, query.Pagination.PageSize)
	assert.Equal(t, prodList.SortByCreatedAt, query.SortBy)
	assert.Equal(t, prodList.SortOrderDesc, query.SortOrder)
}

func TestNewQuery_TrimsSearch(t *testing.T) {
	t.Parallel()

	query := prodList.NewQuery(
		1,
		10,
		"  iphone  ",
		"",
		"",
	)

	assert.Equal(t, "iphone", query.Search)
}

func TestNewQuery_NormalizesSortOrder(t *testing.T) {
	t.Parallel()

	query := prodList.NewQuery(
		1,
		10,
		"",
		prodList.SortByPrice,
		"DESC",
	)

	assert.Equal(
		t,
		prodList.SortOrderDesc,
		query.SortOrder,
	)
}

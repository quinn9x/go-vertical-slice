package list_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	prodList "github.com/quinn9x/go-vertical-slice/internal/products/features/list"
)

func TestNewQuery_UsesDefaultSorting(t *testing.T) {
	t.Parallel()

	query := prodList.NewQuery(
		1,
		10,
		"",
		"",
		"",
	)

	assert.Equal(t, "createdAt", query.SortBy)
	assert.Equal(t, "desc", query.SortOrder)
}

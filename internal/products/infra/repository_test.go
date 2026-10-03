package infra_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	"github.com/quinn9x/go-vertical-slice/internal/products/infra"
	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
)

const iPhone17 = "iPhone 17"

const iPhone17Pro = "iPhone 17 Pro"

func TestProductRepository_Update(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)

	repository := infra.NewProductRepository(db)

	product := &domain.Product{
		ID:      uuid.NewString(),
		Name:    iPhone17,
		Price:   999,
		Version: 1,
	}

	require.NoError(
		t,
		db.Create(product).Error,
	)

	product.Name = iPhone17Pro
	product.Price = 1299

	err := repository.Update(
		context.Background(),
		product,
		1,
	)

	require.NoError(t, err)

	var actual domain.Product

	require.NoError(
		t,
		db.First(&actual, "id = ?", product.ID).Error,
	)

	assert.Equal(t, iPhone17Pro, actual.Name)
	assert.Equal(t, 1299.0, actual.Price)
	assert.Equal(t, int64(2), actual.Version)
}

func TestProductRepository_Update_ReturnsConflictForStaleVersion(
	t *testing.T,
) {
	t.Parallel()

	db := newTestDB(t)

	repository := infra.NewProductRepository(db)

	product := &domain.Product{
		ID:      uuid.NewString(),
		Name:    iPhone17,
		Price:   999,
		Version: 1,
	}

	require.NoError(
		t,
		db.Create(product).Error,
	)

	// First update: version 1 -> 2.
	product.Name = iPhone17Pro

	require.NoError(
		t,
		repository.Update(
			context.Background(),
			product,
			1,
		),
	)

	// Second update still uses stale version 1.
	product.Name = "Another name"

	err := repository.Update(
		context.Background(),
		product,
		1,
	)

	require.Error(t, err)

	var appErr *apperrors.Error

	require.ErrorAs(t, err, &appErr)

	assert.Equal(
		t,
		apperrors.CodeConflict,
		appErr.Code,
	)
}

func TestProductRepository_Update_ReturnsConflictWhenProductDoesNotExist(
	t *testing.T,
) {
	t.Parallel()

	db := newTestDB(t)

	repository := infra.NewProductRepository(db)

	product := &domain.Product{
		ID:      uuid.NewString(),
		Name:    "Missing",
		Price:   100,
		Version: 1,
	}

	err := repository.Update(
		context.Background(),
		product,
		1,
	)

	require.Error(t, err)
}

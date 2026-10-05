package infra_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/quinn9x/go-vertical-slice/internal/products"
	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	"github.com/quinn9x/go-vertical-slice/internal/products/infra"
	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
	"github.com/quinn9x/go-vertical-slice/internal/shared/pagination"
)

const (
	iPhone17         = "iPhone 17"
	iPhone17Pro      = "iPhone 17 Pro"
	sortByPrice      = "price"
	defaultSortOrder = "asc"
)

func TestProductRepository_Create(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	repository := infra.NewProductRepository(db)

	product := &domain.Product{
		ID:      uuid.NewString(),
		Name:    iPhone17,
		Price:   999,
		Version: 1,
	}

	err := repository.Create(
		context.Background(),
		product,
	)

	require.NoError(t, err)

	var actual domain.Product

	require.NoError(
		t,
		db.First(&actual, "id = ?", product.ID).Error,
	)

	assert.Equal(t, product.ID, actual.ID)
	assert.Equal(t, product.Name, actual.Name)
	assert.Equal(t, product.Price, actual.Price)
	assert.Equal(t, int64(1), actual.Version)
}

func TestProductRepository_GetByID(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	repository := infra.NewProductRepository(db)

	product := &domain.Product{
		ID:      uuid.NewString(),
		Name:    iPhone17,
		Price:   999,
		Version: 1,
	}

	require.NoError(t, db.Create(product).Error)

	actual, err := repository.GetByID(
		context.Background(),
		product.ID,
	)

	require.NoError(t, err)
	require.NotNil(t, actual)

	assert.Equal(t, product.ID, actual.ID)
	assert.Equal(t, product.Name, actual.Name)
	assert.Equal(t, product.Price, actual.Price)
	assert.Equal(t, product.Version, actual.Version)
}

func TestProductRepository_GetByID_ReturnsNotFound(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	repository := infra.NewProductRepository(db)

	_, err := repository.GetByID(
		context.Background(),
		uuid.NewString(),
	)

	require.Error(t, err)

	var appErr *apperrors.Error

	require.ErrorAs(t, err, &appErr)

	assert.Equal(
		t,
		apperrors.CodeNotFound,
		appErr.Code,
	)
}

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

func TestProductRepository_List_Pagination(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	repository := infra.NewProductRepository(db)

	seedProducts(t, db)

	items, total, err := repository.List(
		context.Background(),
		products.ListOptions{
			Pagination: pagination.Params{
				Page:     1,
				PageSize: 2,
			},
		},
	)

	require.NoError(t, err)

	assert.Equal(t, int64(4), total)
	assert.Len(t, items, 2)
}

func TestProductRepository_List_Search(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	repository := infra.NewProductRepository(db)

	seedProducts(t, db)

	items, total, err := repository.List(
		context.Background(),
		products.ListOptions{
			Pagination: pagination.Params{
				Page:     1,
				PageSize: 10,
			},
			Search: "Mac",
		},
	)

	require.NoError(t, err)

	assert.Equal(t, int64(2), total)
	assert.Len(t, items, 2)
}

func TestProductRepository_List_SortByPriceAscending(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	repository := infra.NewProductRepository(db)

	seedProducts(t, db)

	items, _, err := repository.List(
		context.Background(),
		products.ListOptions{
			Pagination: pagination.Params{
				Page:     1,
				PageSize: 10,
			},
			SortBy:    sortByPrice,
			SortOrder: defaultSortOrder,
		},
	)

	require.NoError(t, err)

	assertProductsSortedByPrice(t, items, defaultSortOrder)
}

func TestProductRepository_List_SortByPriceDescending(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	repository := infra.NewProductRepository(db)

	seedProducts(t, db)

	items, _, err := repository.List(
		context.Background(),
		products.ListOptions{
			Pagination: pagination.Params{
				Page:     1,
				PageSize: 10,
			},
			SortBy:    sortByPrice,
			SortOrder: "desc",
		},
	)

	require.NoError(t, err)

	assertProductsSortedByPrice(t, items, "desc")
}

func TestProductRepository_Delete(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	repository := infra.NewProductRepository(db)

	product := &domain.Product{
		ID:      uuid.NewString(),
		Name:    iPhone17,
		Price:   999,
		Version: 1,
	}

	require.NoError(t, db.Create(product).Error)

	err := repository.Delete(
		context.Background(),
		product.ID,
	)

	require.NoError(t, err)

	var actual domain.Product

	err = db.First(&actual, "id = ?", product.ID).Error

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestProductRepository_Delete_ReturnsNotFound(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	repository := infra.NewProductRepository(db)

	err := repository.Delete(
		context.Background(),
		uuid.NewString(),
	)

	require.Error(t, err)

	var appErr *apperrors.Error

	require.ErrorAs(t, err, &appErr)

	assert.Equal(
		t,
		apperrors.CodeNotFound,
		appErr.Code,
	)
}

func TestProductRepository_GetByID_ContextCancelled(t *testing.T) {
	t.Parallel()

	db := newTestDB(t)
	repository := infra.NewProductRepository(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repository.GetByID(
		ctx,
		uuid.NewString(),
	)

	require.Error(t, err)
}

func seedProducts(
	t *testing.T,
	db *gorm.DB,
) {
	t.Helper()

	items := []domain.Product{
		{
			ID:      uuid.NewString(),
			Name:    "MacBook Pro",
			Price:   2999,
			Version: 1,
		},
		{
			ID:      uuid.NewString(),
			Name:    iPhone17,
			Price:   999,
			Version: 1,
		},
		{
			ID:      uuid.NewString(),
			Name:    "iPad Pro",
			Price:   1299,
			Version: 1,
		},
		{
			ID:      uuid.NewString(),
			Name:    "Mac Mini",
			Price:   799,
			Version: 1,
		},
	}

	require.NoError(t, db.Create(&items).Error)
}

func assertProductsSortedByPrice(
	t *testing.T,
	items []domain.Product,
	sortOrder string,
) {
	t.Helper()

	require.Len(t, items, 4)

	for i := 1; i < len(items); i++ {
		if sortOrder == defaultSortOrder {
			assert.LessOrEqual(
				t,
				items[i-1].Price,
				items[i].Price,
			)

			continue
		}

		assert.GreaterOrEqual(
			t,
			items[i-1].Price,
			items[i].Price,
		)
	}
}

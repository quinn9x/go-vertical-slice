package list_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/quinn9x/go-vertical-slice/internal/products"
	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	prodList "github.com/quinn9x/go-vertical-slice/internal/products/features/list"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
	"github.com/quinn9x/go-vertical-slice/internal/shared/pagination"
)

var errDatabase = errors.New("database error")

const search = "iphone"

const productID1 = "product-1"

const productID2 = "product-2"

const productName = "iPhone 17"

func TestHandler_Handle_ReturnsProducts(t *testing.T) {
	t.Parallel()

	input := struct {
		Pagination pagination.Params
		Search     string
		SortBy     string
		SortOrder  string
	}{
		Pagination: pagination.New(1, 10),
		Search:     search,
		SortBy:     "price",
		SortOrder:  "asc",
	}

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)

	handler := prodList.NewHandler(repository)

	options := products.ListOptions{
		Pagination: input.Pagination,
		Search:     input.Search,
		SortBy:     input.SortBy,
		SortOrder:  input.SortOrder,
	}

	total := int64(2)

	repository.
		EXPECT().
		List(gomock.Any(), options).
		Return([]domain.Product{
			{
				ID:    productID1,
				Name:  productName,
				Price: 999,
			},
			{
				ID:    productID2,
				Name:  "iPhone 17 Pro",
				Price: 1299,
			},
		}, total, nil)

	result, err := handler.Handle(
		context.Background(),
		prodList.ListerProductsQuery{
			Pagination: input.Pagination,
			Search:     input.Search,
			SortBy:     input.SortBy,
			SortOrder:  input.SortOrder,
		},
	)

	require.NoError(t, err)

	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)

	require.Len(t, result.Items, 2)

	assert.Equal(t, "iPhone 17", result.Items[0].Name)
	assert.Equal(t, 999.0, result.Items[0].Price)
}

func TestHandler_Handle_UsesPagination(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	handler := prodList.NewHandler(repository)

	options := products.ListOptions{
		Pagination: pagination.New(2, 10),
		SortBy:     prodList.SortByCreatedAt,
		SortOrder:  prodList.SortOrderDesc,
	}

	total := int64(11)

	repository.
		EXPECT().
		List(gomock.Any(), options).
		Return(
			[]domain.Product{
				{
					ID:    "product-11",
					Name:  "Product 11",
					Price: 110,
				},
			},
			total,
			nil,
		)

	result, err := handler.Handle(
		context.Background(),
		prodList.ListerProductsQuery{
			Pagination: pagination.New(2, 10),
			SortBy:     prodList.SortByCreatedAt,
			SortOrder:  prodList.SortOrderDesc,
		},
	)

	require.NoError(t, err)

	assert.Equal(t, 2, result.Page)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, int64(11), result.Total)
	assert.Len(t, result.Items, 1)
}

func TestHandler_Handle_UsesSearch(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	handler := prodList.NewHandler(repository)

	options := products.ListOptions{
		Pagination: pagination.New(1, 10),
		Search:     search,
		SortBy:     prodList.SortByCreatedAt,
		SortOrder:  prodList.SortOrderDesc,
	}

	total := int64(1)

	repository.
		EXPECT().
		List(gomock.Any(), options).
		Return(
			[]domain.Product{
				{
					ID:    productID1,
					Name:  productName,
					Price: 999,
				},
			},
			total,
			nil,
		)

	result, err := handler.Handle(
		context.Background(),
		prodList.ListerProductsQuery{
			Pagination: pagination.New(1, 10),
			Search:     search,
			SortBy:     prodList.SortByCreatedAt,
			SortOrder:  prodList.SortOrderDesc,
		},
	)

	require.NoError(t, err)

	require.Len(t, result.Items, 1)
	assert.Equal(t, productName, result.Items[0].Name)
	assert.Equal(t, int64(1), result.Total)
}

func TestHandler_Handle_UsesSorting(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	handler := prodList.NewHandler(repository)

	options := products.ListOptions{
		Pagination: pagination.New(1, 10),
		SortBy:     prodList.SortByPrice,
		SortOrder:  prodList.SortOrderAsc,
	}

	total := int64(2)

	repository.
		EXPECT().
		List(gomock.Any(), options).
		Return(
			[]domain.Product{
				{
					ID:    productID1,
					Name:  "Cheap Product",
					Price: 100,
				},
				{
					ID:    "product-2",
					Name:  "Expensive Product",
					Price: 1000,
				},
			},
			total,
			nil,
		)

	result, err := handler.Handle(
		context.Background(),
		prodList.ListerProductsQuery{
			Pagination: pagination.New(1, 10),
			SortBy:     prodList.SortByPrice,
			SortOrder:  prodList.SortOrderAsc,
		},
	)

	require.NoError(t, err)

	require.Len(t, result.Items, 2)

	assert.Equal(t, 100.0, result.Items[0].Price)
	assert.Equal(t, 1000.0, result.Items[1].Price)
}

func TestHandler_Handle_ReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	handler := prodList.NewHandler(repository)

	options := products.ListOptions{
		Pagination: pagination.New(1, 10),
		SortBy:     prodList.SortByCreatedAt,
		SortOrder:  prodList.SortOrderDesc,
	}

	total := int64(0)

	repository.
		EXPECT().
		List(gomock.Any(), options).
		Return(nil, total, errDatabase)

	_, err := handler.Handle(
		context.Background(),
		prodList.ListerProductsQuery{
			Pagination: pagination.New(1, 10),
			SortBy:     prodList.SortByCreatedAt,
			SortOrder:  prodList.SortOrderDesc,
		},
	)

	require.ErrorIs(t, err, errDatabase)
}

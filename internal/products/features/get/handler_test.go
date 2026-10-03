package get_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	prodGetter "github.com/quinn9x/go-vertical-slice/internal/products/features/get"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
)

var errProductNotFound = errors.New("product not found")

func TestHandler_Handle_ReturnsProduct(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repository := mocks.NewMockRepository(ctrl)

	handler := prodGetter.NewHandler(repository)

	const productID = "product-1"

	product := &domain.Product{
		ID:    productID,
		Name:  "iPhone 17",
		Price: 999,
	}

	repository.
		EXPECT().
		GetByID(gomock.Any(), productID).
		Return(product, nil)

	result, err := handler.Handle(
		context.Background(),
		prodGetter.GetterProductQuery{
			ID: productID,
		},
	)

	require.NoError(t, err)

	assert.Equal(t, productID, result.ID)
	assert.Equal(t, "iPhone 17", result.Name)
	assert.Equal(t, 999.0, result.Price)
}

func TestHandler_Handle_ReturnsErrorWhenProductNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)

	handler := prodGetter.NewHandler(repository)

	repository.
		EXPECT().
		GetByID(gomock.Any(), "missing").
		Return(nil, errProductNotFound)

	_, err := handler.Handle(
		context.Background(),
		prodGetter.GetterProductQuery{
			ID: "missing",
		},
	)

	require.ErrorIs(t, err, errProductNotFound)
}

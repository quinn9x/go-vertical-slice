package update_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	prodUpdate "github.com/quinn9x/go-vertical-slice/internal/products/features/update"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

var errRequestValidation = errors.New("request validation failed")

var errDatabase = errors.New("database error")

const productID = "product-1"

const iPhone17 = "iPhone 17"

const iPhone17Pro = "iPhone 17 Pro"

func TestHandler_Handle_UpdatesProduct(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	handler := prodUpdate.NewHandler(validate, repository)

	product := &domain.Product{
		ID:      productID,
		Name:    iPhone17,
		Price:   999,
		Version: 1,
	}

	repository.
		EXPECT().
		GetByID(gomock.Any(), productID).
		Return(product, nil)

	repository.
		EXPECT().
		Update(gomock.Any(), product, int64(1)).
		Do(func(
			_ context.Context,
			updated *domain.Product,
			expectedVersion int64,
		) {
			assert.Equal(t, iPhone17Pro, updated.Name)
			assert.Equal(t, 1299.0, updated.Price)
			assert.Equal(t, int64(1), expectedVersion)
		}).
		Return(nil)

	result, err := handler.Handle(
		context.Background(),
		prodUpdate.UpdaterProductCommand{
			ID:      productID,
			Name:    iPhone17Pro,
			Price:   1299,
			Version: 1,
		},
	)

	require.NoError(t, err)

	assert.Equal(t, productID, result.ID)
	assert.Equal(t, iPhone17Pro, result.Name)
	assert.Equal(t, 1299.0, result.Price)
	assert.Equal(t, int64(2), result.Version)
}

func TestHandler_Handle_ReturnsErrorWhenGetProductFails(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	handler := prodUpdate.NewHandler(validate, repository)

	repository.
		EXPECT().
		GetByID(gomock.Any(), productID).
		Return(nil, errDatabase)

	_, err := handler.Handle(
		context.Background(),
		prodUpdate.UpdaterProductCommand{
			ID:      productID,
			Name:    iPhone17Pro,
			Price:   1299,
			Version: 1,
		},
	)

	require.ErrorIs(t, err, errDatabase)
}

func TestHandler_Handle_ReturnsConflict(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	handler := prodUpdate.NewHandler(
		validate,
		repository,
	)

	product := &domain.Product{
		ID:      productID,
		Name:    iPhone17,
		Price:   999,
		Version: 2,
	}

	repository.
		EXPECT().
		GetByID(
			gomock.Any(),
			productID,
		).
		Return(product, nil)

	expectedErr := apperrors.NewConflict(
		"Product was modified by another request",
	)

	repository.
		EXPECT().
		Update(
			gomock.Any(),
			product,
			int64(1),
		).
		Return(expectedErr)

	_, err := handler.Handle(
		context.Background(),
		prodUpdate.UpdaterProductCommand{
			ID:      productID,
			Name:    iPhone17Pro,
			Price:   1299,
			Version: 1,
		},
	)

	require.Error(t, err)
}

func TestHandler_Handle_RejectsInvalidVersion(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	handler := prodUpdate.NewHandler(
		validate,
		repository,
	)

	_, err := handler.Handle(
		context.Background(),
		prodUpdate.UpdaterProductCommand{
			ID:      productID,
			Name:    iPhone17Pro,
			Price:   1299,
			Version: 0,
		},
	)

	require.Error(t, err)
}

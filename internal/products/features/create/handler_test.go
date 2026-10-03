package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	prodCreate "github.com/quinn9x/go-vertical-slice/internal/products/features/create"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

var errDatabase = errors.New("database error")

const productName = "iPhone 17"

func TestHandler_Handle_CreatesProduct(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	handler := prodCreate.NewHandler(
		validate,
		repository,
	)

	repository.
		EXPECT().
		Create(
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				_ context.Context,
				product *domain.Product,
			) error {
				assert.Equal(t, productName, product.Name)
				assert.Equal(t, 999.0, product.Price)

				return nil
			},
		)

	result, err := handler.Handle(
		context.Background(),
		prodCreate.CreatorProductCommand{
			Name:  productName,
			Price: 999,
		},
	)

	require.NoError(t, err)

	assert.Equal(t, productName, result.Name)
	assert.Equal(t, 999.0, result.Price)
}

func TestHandler_Handle_RejectsInvalidCommand(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	handler := prodCreate.NewHandler(
		validate,
		repository,
	)

	_, err := handler.Handle(
		context.Background(),
		prodCreate.CreatorProductCommand{
			Name:  "",
			Price: 0,
		},
	)

	require.Error(t, err)
}

func TestHandler_Handle_ReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	handler := prodCreate.NewHandler(
		validate,
		repository,
	)

	repository.
		EXPECT().
		Create(
			gomock.Any(),
			gomock.Any(),
		).
		Return(errDatabase)

	_, err := handler.Handle(
		context.Background(),
		prodCreate.CreatorProductCommand{
			Name:  productName,
			Price: 999,
		},
	)

	require.ErrorIs(t, err, errDatabase)
}

package deleter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	prodDeleter "github.com/quinn9x/go-vertical-slice/internal/products/features/delete"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
)

var errProductNotFound = errors.New("product not found")

func TestHandler_Handle_DeletesProduct(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)

	handler := prodDeleter.NewHandler(repository)

	repository.
		EXPECT().
		Delete(gomock.Any(), "product-1").
		Return(nil)

	err := handler.Handle(
		context.Background(),
		prodDeleter.DeleteProductCommand{
			ID: "product-1",
		},
	)

	require.NoError(t, err)
}

func TestHandler_Handle_ReturnsErrorWhenProductNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)

	handler := prodDeleter.NewHandler(repository)

	repository.
		EXPECT().
		Delete(gomock.Any(), "missing").
		Return(errProductNotFound)

	err := handler.Handle(
		context.Background(),
		prodDeleter.DeleteProductCommand{
			ID: "missing",
		},
	)

	require.ErrorIs(t, err, errProductNotFound)
}

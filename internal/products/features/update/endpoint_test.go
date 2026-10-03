package update_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	prodUpdate "github.com/quinn9x/go-vertical-slice/internal/products/features/update"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

func TestEndpoint_Handle_UpdatesProduct(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	product := &domain.Product{
		ID:      "123",
		Name:    "iPhone 16",
		Price:   999,
		Version: 1,
	}

	repository.
		EXPECT().
		GetByID(
			gomock.Any(),
			"123",
		).
		Return(product, nil)

	repository.
		EXPECT().
		Update(
			gomock.Any(),
			product,
			int64(1),
		).
		Return(nil)

	handler := prodUpdate.NewHandler(
		validate,
		repository,
	)

	endpoint := prodUpdate.NewEndpoint(handler)

	e := echo.New()
	e.PUT("/api/products/:id", endpoint.Handle)

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPut,
		"/api/products/123",
		strings.NewReader(`{
			"name": "iPhone 17 Pro",
			"price": 1299,
			"version": 1
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

package update_test

import (
	"context"
	"encoding/json"
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
	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

const iPhone16 = "iPhone 16"

func TestEndpoint_Handle_UpdatesProduct(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	product := &domain.Product{
		ID:      "123",
		Name:    iPhone16,
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

func TestEndpoint_Handle_ReturnsConflict(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	product := &domain.Product{
		ID:      "123",
		Name:    iPhone16,
		Price:   999,
		Version: 2,
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
		Return(
			apperrors.NewConflict(
				"Product was modified by another request",
			),
		)

	handler := prodUpdate.NewHandler(
		validate,
		repository,
	)

	endpoint := prodUpdate.NewEndpoint(handler)

	e := echo.New()
	e.HTTPErrorHandler = apperrors.HTTPErrorHandler

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

	require.Equal(t, http.StatusConflict, rec.Code)

	var response struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	require.NoError(
		t,
		json.Unmarshal(rec.Body.Bytes(), &response),
	)

	require.Equal(
		t,
		apperrors.CodeConflict,
		response.Error.Code,
	)

	require.Equal(
		t,
		"Product was modified by another request",
		response.Error.Message,
	)
}

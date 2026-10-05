package app_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	prodCreate "github.com/quinn9x/go-vertical-slice/internal/products/features/create"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

func TestCreateProduct_InvalidRequest_ReturnsValidationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validator := validation.New()

	handler := prodCreate.NewHandler(
		validator,
		repository,
	)

	endpoint := prodCreate.NewEndpoint(handler)

	e := echo.New()
	e.HTTPErrorHandler = apperrors.HTTPErrorHandler

	e.POST("/api/products", endpoint.Handle)

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/products",
		strings.NewReader(`{
			"name": "",
			"price": 0
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var response struct {
		Error struct {
			Code    string            `json:"code"`
			Message string            `json:"message"`
			Fields  map[string]string `json:"fields"`
		} `json:"error"`
	}

	require.NoError(
		t,
		json.Unmarshal(rec.Body.Bytes(), &response),
	)

	assert.Equal(
		t,
		apperrors.CodeValidation,
		response.Error.Code,
	)

	assert.Equal(
		t,
		"Request validation failed",
		response.Error.Message,
	)

	assert.Contains(t, response.Error.Fields, "name")
	assert.Contains(t, response.Error.Fields, "price")
}

func TestCreateProduct_MalformedJSON_ReturnsBadRequest(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validator := validation.New()

	handler := prodCreate.NewHandler(
		validator,
		repository,
	)

	endpoint := prodCreate.NewEndpoint(handler)

	e := echo.New()
	e.HTTPErrorHandler = apperrors.HTTPErrorHandler

	e.POST("/api/products", endpoint.Handle)

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/products",
		strings.NewReader(`{
			"name": "iPhone 17",
			"price": 999`),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var response struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	require.NoError(
		t,
		json.Unmarshal(rec.Body.Bytes(), &response),
	)

	assert.Equal(
		t,
		apperrors.CodeBadRequest,
		response.Error.Code,
	)
}

package apperrors_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
)

var errDatabase = errors.New("database password=secret host=postgres")

func TestHTTPErrorHandler_ValidationError(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/products", http.NoBody)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := apperrors.NewValidationError(
		map[string]string{
			"name":  "is required",
			"price": "must be greater than 0",
		},
	)

	apperrors.HTTPErrorHandler(ctx, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response struct {
		Error apperrors.Error `json:"error"`
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

	assert.Equal(
		t,
		"is required",
		response.Error.Fields["name"],
	)

	assert.Equal(
		t,
		"must be greater than 0",
		response.Error.Fields["price"],
	)
}

func TestHTTPErrorHandler_NotFound(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/products/1", http.NoBody)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := apperrors.NewNotFound("product not found")

	apperrors.HTTPErrorHandler(ctx, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response struct {
		Error apperrors.Error `json:"error"`
	}

	require.NoError(
		t,
		json.Unmarshal(rec.Body.Bytes(), &response),
	)

	assert.Equal(
		t,
		apperrors.CodeNotFound,
		response.Error.Code,
	)

	assert.Equal(
		t,
		"product not found",
		response.Error.Message,
	)
}

func TestHTTPErrorHandler_Conflict(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/products/1", http.NoBody)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := apperrors.NewConflict(
		"Product was modified by another request",
	)

	apperrors.HTTPErrorHandler(ctx, err)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var response struct {
		Error apperrors.Error `json:"error"`
	}

	require.NoError(
		t,
		json.Unmarshal(rec.Body.Bytes(), &response),
	)

	assert.Equal(
		t,
		apperrors.CodeConflict,
		response.Error.Code,
	)
}

func TestHTTPErrorHandler_InternalError(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	apperrors.HTTPErrorHandler(ctx, errDatabase)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response struct {
		Error apperrors.Error `json:"error"`
	}

	require.NoError(
		t,
		json.Unmarshal(rec.Body.Bytes(), &response),
	)

	assert.Equal(
		t,
		apperrors.CodeInternalError,
		response.Error.Code,
	)

	assert.Equal(
		t,
		"Internal server error",
		response.Error.Message,
	)

	assert.NotContains(
		t,
		rec.Body.String(),
		"password",
	)
}

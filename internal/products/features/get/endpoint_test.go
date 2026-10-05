package get_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	prodGetter "github.com/quinn9x/go-vertical-slice/internal/products/features/get"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
)

const productID = "product-1"

const keyID = "id"

func TestEndpoint_Handle_ReturnsOK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)

	repository.
		EXPECT().
		GetByID(gomock.Any(), productID).
		Return(&domain.Product{
			ID: productID,
		}, nil)

	handler := prodGetter.NewHandler(repository)
	endpoint := prodGetter.NewEndpoint(handler)

	e := echo.New()

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/products/"+productID,
		http.NoBody,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/api/products/:id")
	c.SetPathValues(echo.PathValues{
		{
			Name:  keyID,
			Value: productID,
		},
	})

	err := endpoint.Handle(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestEndpoint_Handle_ReturnsError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)

	repository.
		EXPECT().
		GetByID(gomock.Any(), productID).
		Return(nil, errProductNotFound)

	handler := prodGetter.NewHandler(repository)
	endpoint := prodGetter.NewEndpoint(handler)

	e := echo.New()

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/products/"+productID,
		http.NoBody,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/api/products/:id")
	c.SetPathValues(echo.PathValues{
		{
			Name:  keyID,
			Value: productID,
		},
	})

	err := endpoint.Handle(c)

	require.Error(t, err)
}

func TestEndpoint_Handle_ReturnsNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)

	repository.
		EXPECT().
		GetByID(gomock.Any(), productID).
		Return(
			nil,
			apperrors.NewNotFound("Product not found"),
		)

	handler := prodGetter.NewHandler(repository)
	endpoint := prodGetter.NewEndpoint(handler)

	e := echo.New()
	e.HTTPErrorHandler = apperrors.HTTPErrorHandler

	e.GET("/api/products/:id", endpoint.Handle)

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/products/"+productID,
		http.NoBody,
	)

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

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
		apperrors.CodeNotFound,
		response.Error.Code,
	)

	require.Equal(
		t,
		"Product not found",
		response.Error.Message,
	)
}

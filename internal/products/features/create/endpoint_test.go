package create_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	prodCreate "github.com/quinn9x/go-vertical-slice/internal/products/features/create"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

func TestEndpoint_Handle_CreatesProduct(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	validate := validation.New()

	repository.
		EXPECT().
		Create(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil)

	handler := prodCreate.NewHandler(
		validate,
		repository,
	)

	endpoint := prodCreate.NewEndpoint(handler)

	e := echo.New()

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/products",
		strings.NewReader(`{
			"name": "iPhone 17",
			"price": 999
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := endpoint.Handle(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
}

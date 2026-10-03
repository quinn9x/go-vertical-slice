package list_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	prodList "github.com/quinn9x/go-vertical-slice/internal/products/features/list"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
)

func TestEndpoint_Handle_RejectsInvalidSortBy(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	handler := prodList.NewHandler(repository)
	endpoint := prodList.NewEndpoint(handler)

	e := echo.New()

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/products?sortBy=invalid",
		http.NoBody,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := endpoint.Handle(c)

	require.Error(t, err)
}

func TestEndpoint_Handle_RejectsInvalidSortOrder(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	handler := prodList.NewHandler(repository)
	endpoint := prodList.NewEndpoint(handler)

	e := echo.New()

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/products?sortOrder=random",
		http.NoBody,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := endpoint.Handle(c)

	require.Error(t, err)
}

package deleter_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	prodDeleter "github.com/quinn9x/go-vertical-slice/internal/products/features/delete"
	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
)

const productID = "product-1"

func TestEndpoint_Handle_ReturnsNoContent(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)

	repository.
		EXPECT().
		Delete(gomock.Any(), productID).
		Return(nil)

	handler := prodDeleter.NewHandler(repository)
	endpoint := prodDeleter.NewEndpoint(handler)

	e := echo.New()

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		"/api/products/"+productID,
		http.NoBody,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/api/products/:id")
	c.SetPathValues(echo.PathValues{
		{
			Name:  "id",
			Value: productID,
		},
	})

	err := endpoint.Handle(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

package list

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/quinn9x/go-vertical-slice/internal/products/mocks"
)

func TestParseQueryInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{
			name:  "empty",
			value: "",
			want:  0,
		},
		{
			name:  "valid",
			value: "10",
			want:  10,
		},
		{
			name:    "invalid",
			value:   "abc",
			wantErr: true,
		},
		{
			name:  "negative",
			value: "-1",
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseQueryInt(tt.value)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEndpoint_Handle_RejectsInvalidPage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	handler := NewHandler(repository)
	endpoint := NewEndpoint(handler)

	e := echo.New()

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/products?page=abc",
		http.NoBody,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := endpoint.Handle(c)

	require.Error(t, err)
}

func TestEndpoint_Handle_RejectsInvalidPageSize(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	handler := NewHandler(repository)
	endpoint := NewEndpoint(handler)

	e := echo.New()

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/products?pageSize=abc",
		http.NoBody,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := endpoint.Handle(c)

	require.Error(t, err)
}

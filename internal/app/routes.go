package app

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/quinn9x/go-vertical-slice/internal/products"
	prodCreate "github.com/quinn9x/go-vertical-slice/internal/products/features/create"
	prodDelete "github.com/quinn9x/go-vertical-slice/internal/products/features/delete"
	prodGet "github.com/quinn9x/go-vertical-slice/internal/products/features/get"
	prodList "github.com/quinn9x/go-vertical-slice/internal/products/features/list"
	prodUpdate "github.com/quinn9x/go-vertical-slice/internal/products/features/update"
	"github.com/quinn9x/go-vertical-slice/internal/shared/database"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

const keyStatus = "status"

func registerProductRoutes(
	e *echo.Echo,
	validator *validation.Validator,
	repository products.Repository,
) {
	createHandler := prodCreate.NewHandler(
		validator,
		repository,
	)
	createEndpoint := prodCreate.NewEndpoint(createHandler)

	getHandler := prodGet.NewHandler(repository)
	getEndpoint := prodGet.NewEndpoint(getHandler)

	listHandler := prodList.NewHandler(repository)
	listEndpoint := prodList.NewEndpoint(listHandler)

	updateHandler := prodUpdate.NewHandler(
		validator,
		repository,
	)
	updateEndpoint := prodUpdate.NewEndpoint(updateHandler)

	deleteHandler := prodDelete.NewHandler(repository)
	deleteEndpoint := prodDelete.NewEndpoint(deleteHandler)

	e.POST("/api/products", createEndpoint.Handle)
	e.GET("/api/products/:id", getEndpoint.Handle)
	e.GET("/api/products", listEndpoint.Handle)
	e.PUT("/api/products/:id", updateEndpoint.Handle)
	e.DELETE("/api/products/:id", deleteEndpoint.Handle)
}

func registerHealthRoute(e *echo.Echo, db *database.DB) {
	e.GET("/health", func(c *echo.Context) error {
		if err := db.Ping(); err != nil {
			return c.JSON(
				http.StatusServiceUnavailable,
				map[string]string{
					keyStatus: "unhealthy",
				},
			)
		}

		return c.JSON(
			http.StatusOK,
			map[string]string{
				keyStatus: "ok",
			},
		)
	})
}

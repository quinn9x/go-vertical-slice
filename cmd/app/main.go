// Package main provides the application entry point for a Go vertical slice architecture template.
package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"

	prodCreate "github.com/quinn9x/go-vertical-slice/internal/products/features/create"
	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

func main() {
	e := echo.New()

	e.HTTPErrorHandler = apperrors.HTTPErrorHandler

	validator := validation.New()

	prodCreateHandler := prodCreate.NewHandler(validator)
	prodCreateEndpoint := prodCreate.NewEndpoint(prodCreateHandler)

	e.POST("/api/products", prodCreateEndpoint.Handle)

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	if err := e.Start(":9080"); err != nil {
		log.Fatal(err)
	}
}

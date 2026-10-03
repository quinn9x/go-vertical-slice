// Package main provides the application entry point for a Go vertical slice architecture template.
package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
	prodCreate "github.com/quinn9x/go-vertical-slice/internal/products/features/create"
	prodGet "github.com/quinn9x/go-vertical-slice/internal/products/features/get"
	"github.com/quinn9x/go-vertical-slice/internal/products/infra"
	"github.com/quinn9x/go-vertical-slice/internal/shared/database"
	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

func main() {
	e := echo.New()

	e.HTTPErrorHandler = apperrors.HTTPErrorHandler

	db, err := database.New("data/app.db")
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&domain.Product{}); err != nil {
		log.Fatal(err)
	}

	prodRepo := infra.NewProductRepository(db.DB)

	validator := validation.New()

	prodCreateHandler := prodCreate.NewHandler(validator, prodRepo)
	prodCreateEndpoint := prodCreate.NewEndpoint(prodCreateHandler)

	prodGetHandler := prodGet.NewHandler(prodRepo)
	prodGetEndpoint := prodGet.NewEndpoint(prodGetHandler)

	e.POST("/api/products", prodCreateEndpoint.Handle)
	e.GET("/api/products/:id", prodGetEndpoint.Handle)

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	if err := e.Start(":9080"); err != nil {
		log.Fatal(err)
	}
}

// Package main provides the application entry point for a Go vertical slice architecture template.
package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	if err := e.Start(":9080"); err != nil {
		log.Fatal(err)
	}
}

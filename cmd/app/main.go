// Package main provides the application entry point for a Go vertical slice architecture template.
package main

import (
	"fmt"
	"log"

	"github.com/quinn9x/go-vertical-slice/internal/app"
	"github.com/quinn9x/go-vertical-slice/internal/shared/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	application, err := app.New(&cfg)
	if err != nil {
		return fmt.Errorf("create application: %w", err)
	}

	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("failed to close application: %v", err)
		}
	}()

	if err := application.Start(); err != nil {
		return fmt.Errorf("start application: %w", err)
	}

	return nil
}

// Package main provides the application entry point for a Go vertical slice architecture template.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/quinn9x/go-vertical-slice/internal/app"
	"github.com/quinn9x/go-vertical-slice/internal/shared/config"
	"github.com/quinn9x/go-vertical-slice/internal/shared/logger"
)

func main() {
	log := logger.New()

	if err := run(log); err != nil {
		log.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	application, err := app.New(&cfg, log)
	if err != nil {
		return fmt.Errorf("create application: %w", err)
	}

	defer func() {
		if err := application.Close(); err != nil {
			log.Error(
				"failed to close application",
				"error", err,
			)
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := application.Start(ctx); err != nil {
		return fmt.Errorf("start application: %w", err)
	}

	return nil
}

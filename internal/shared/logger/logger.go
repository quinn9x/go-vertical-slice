// Package logger provides application logging utilities.
package logger

import (
	"log/slog"
	"os"
)

// New creates a new application logger.
func New() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)
}

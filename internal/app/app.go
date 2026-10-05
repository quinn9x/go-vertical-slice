// Package app provides application initialization and lifecycle management.
package app

import (
	"fmt"

	"github.com/labstack/echo/v5"

	"github.com/quinn9x/go-vertical-slice/internal/products/infra"
	"github.com/quinn9x/go-vertical-slice/internal/shared/config"
	"github.com/quinn9x/go-vertical-slice/internal/shared/database"
	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
	"github.com/quinn9x/go-vertical-slice/internal/shared/validation"
)

// Application represents the application and its dependencies.
type Application struct {
	Config *config.Config
	DB     *database.DB
	Echo   *echo.Echo
}

// New creates and initializes a new application.
func New(cfg *config.Config) (*Application, error) {
	db, err := database.New(&cfg.Database)
	if err != nil {
		return nil, err
	}

	if err := database.Migrate(db); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("migration failed: %w; close database: %w", err, closeErr)
		}

		return nil, err
	}

	e := echo.New()
	e.HTTPErrorHandler = apperrors.HTTPErrorHandler

	validator := validation.New()

	productRepository := infra.NewProductRepository(db.DB)

	registerProductRoutes(
		e,
		validator,
		productRepository,
	)

	registerHealthRoute(e, db)

	return &Application{
		Config: cfg,
		DB:     db,
		Echo:   e,
	}, nil
}

// Start starts the application server.
func (a *Application) Start() error {
	address := fmt.Sprintf("%s:%d", a.Config.App.Host, a.Config.App.Port)

	return a.Echo.Start(address)
}

// Close closes the application and its database connection.
func (a *Application) Close() error {
	return a.DB.Close()
}

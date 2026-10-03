package infra_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/quinn9x/go-vertical-slice/internal/products/domain"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	ctx := context.Background()

	container, err := postgrescontainer.Run(
		ctx,
		"postgres:18",
		postgrescontainer.WithDatabase("products_test"),
		postgrescontainer.WithUsername("postgres"),
		postgrescontainer.WithPassword("postgres"),
		testcontainers.WithAdditionalWaitStrategy(
			wait.ForLog(
				"database system is ready to accept connections",
			).WithOccurrence(2),
			wait.ForListeningPort("5432/tcp"),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(
			t,
			testcontainers.TerminateContainer(container),
		)
	})

	dsn, err := container.ConnectionString(
		ctx,
		"sslmode=disable",
	)
	require.NoError(t, err)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	require.NoError(t, err)

	require.NoError(
		t,
		db.AutoMigrate(&domain.Product{}),
	)

	return db
}

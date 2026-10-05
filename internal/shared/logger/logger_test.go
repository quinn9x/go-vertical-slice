package logger_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/quinn9x/go-vertical-slice/internal/shared/logger"
)

func TestNew(t *testing.T) {
	t.Parallel()

	log := logger.New()

	require.NotNil(t, log)
}

package endpoint_checkers_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/endpoint_checkers"
	"github.com/flowck/dobermann/backend/internal/common/logs"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func TestHttpChecker(t *testing.T) {
	simulatorEndpointUrl := os.Getenv("SIMULATOR_ENDPOINT_URL")

	httpChecker, err := endpoint_checkers.NewHttpChecker(monitor.RegionEurope.String(), 2, logs.New(false))
	require.NoError(t, err)

	t.Run("is_up", func(t *testing.T) {
		t.Parallel()
		result, err := httpChecker.Check(context.Background(), fmt.Sprintf("%s?is_up=true", simulatorEndpointUrl))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, int(*result.StatusCode()))
		assert.False(t, result.IsEndpointDown())
	})

	t.Run("is_down", func(t *testing.T) {
		t.Parallel()
		result, err := httpChecker.Check(context.Background(), fmt.Sprintf("%s?is_up=false", simulatorEndpointUrl))
		require.NoError(t, err)
		assert.True(t, result.IsEndpointDown())
	})

	t.Run("error_endpoint_timeouts", func(t *testing.T) {
		t.Parallel()

		result, err := httpChecker.Check(context.Background(), fmt.Sprintf("%s?is_up=true&timeout=true", simulatorEndpointUrl))
		require.NoError(t, err)
		assert.True(t, result.IsEndpointDown())
	})

	t.Run("error_endpoint_timeouts", func(t *testing.T) {
		t.Parallel()

		result, err := httpChecker.Check(context.Background(), fmt.Sprintf("%s?is_up=true&timeout=true", simulatorEndpointUrl))
		require.NoError(t, err)
		assert.True(t, result.IsEndpointDown())
	})
}

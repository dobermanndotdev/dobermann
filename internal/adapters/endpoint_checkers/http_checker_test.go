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
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func TestHttpChecker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	simulatorEndpointUrl := os.Getenv("SIMULATOR_ENDPOINT_URL")

	httpChecker, err := endpoint_checkers.NewHttpChecker(monitor.RegionEurope.String(), 2)
	require.NoError(t, err)

	t.Run("is_up", func(t *testing.T) {
		t.Parallel()
		_, err = httpChecker.Check(ctx, fmt.Sprintf("%s?is_up=true", simulatorEndpointUrl))
		assert.NoError(t, err)
	})

	t.Run("is_down", func(t *testing.T) {
		t.Parallel()
		checkResult, err := httpChecker.Check(ctx, fmt.Sprintf("%s?is_up=false", simulatorEndpointUrl))
		require.NoError(t, err)
		assert.True(t, checkResult.IsEndpointDown())
	})

	t.Run("error_endpoint_timeouts", func(t *testing.T) {
		t.Parallel()

		result, err := httpChecker.Check(ctx, fmt.Sprintf("%s?is_up=true&timeout=true", simulatorEndpointUrl))
		require.NoError(t, err)
		assert.Equal(t, http.StatusRequestTimeout, int(result.StatusCode()))
	})
}

package endpoint_checkers_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dobermanndotdev/dobermann/internal/adapters/endpoint_checkers"
	"github.com/dobermanndotdev/dobermann/internal/common/logs"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
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

	/*t.Run("demo", func(t *testing.T) {
		t.Parallel()

		checker, err := endpoint_checkers.NewHttpChecker(monitor.RegionEurope.String(), 30, logs.New(true))
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(time.Second * 2)
			cancel()
		}()

		for {
			result, err := checker.Check(ctx, "https://google.com")
			require.NoError(t, err)
			assert.False(t, result.IsEndpointDown())

			t.Logf("response time --> %dms", result.ResponseTimeInMs())

			time.Sleep(time.Second * 5)
		}
	})*/
}

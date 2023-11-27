package monitors_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/flowck/dobermann/backend/internal/adapters/monitors"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func TestHttpChecker(t *testing.T) {
	simulatorEndpointUrl := os.Getenv("SIMULATOR_ENDPOINT_URL")

	httpChecker := monitors.NewHttpChecker()
	err := httpChecker.Check(ctx, fmt.Sprintf("%s?is_up=true", simulatorEndpointUrl))
	assert.NoError(t, err)

	err = httpChecker.Check(ctx, fmt.Sprintf("%s?is_up=false", simulatorEndpointUrl))
	assert.ErrorIs(t, err, monitor.ErrEndpointIsDown)
}

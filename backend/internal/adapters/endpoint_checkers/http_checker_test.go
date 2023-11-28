package endpoint_checkers_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/flowck/dobermann/backend/internal/adapters/endpoint_checkers"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func TestHttpChecker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	simulatorEndpointUrl := os.Getenv("SIMULATOR_ENDPOINT_URL")

	httpChecker := endpoint_checkers.NewHttpChecker()
	err := httpChecker.Check(ctx, fmt.Sprintf("%s?is_up=true", simulatorEndpointUrl))
	assert.NoError(t, err)

	err = httpChecker.Check(ctx, fmt.Sprintf("%s?is_up=false", simulatorEndpointUrl))
	assert.ErrorIs(t, err, monitor.ErrEndpointIsDown)
}

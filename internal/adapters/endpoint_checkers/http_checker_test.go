package endpoint_checkers_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/endpoint_checkers"
)

func TestHttpChecker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	simulatorEndpointUrl := os.Getenv("SIMULATOR_ENDPOINT_URL")

	httpChecker, err := endpoint_checkers.NewHttpChecker("europe")
	require.NoError(t, err)
	_, err = httpChecker.Check(ctx, fmt.Sprintf("%s?is_up=true", simulatorEndpointUrl))
	assert.NoError(t, err)

	checkResult, err := httpChecker.Check(ctx, fmt.Sprintf("%s?is_up=false", simulatorEndpointUrl))
	require.NoError(t, err)
	assert.True(t, checkResult.IsEndpointDown())
}

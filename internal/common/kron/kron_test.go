package kron_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dobermanndotdev/dobermann/internal/common/kron"
)

func TestKronLifecycle(t *testing.T) {
	counter := 0
	expectedCounterResult := 5
	ctx := context.Background()

	cron := kron.NewService()
	cron.AddJob(kron.NewJob(time.Millisecond*250, func(ctx context.Context) error {
		if counter < expectedCounterResult {
			counter++
		}
		return nil
	}))

	go func() {
		require.NoError(t, cron.Start(ctx))
	}()
	defer func() {
		require.NoError(t, cron.Stop())
	}()

	assert.Eventually(t, func() bool {
		return counter == expectedCounterResult
	}, time.Second*2, time.Millisecond*500)
}

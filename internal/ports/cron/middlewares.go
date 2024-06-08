package cron

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/common/observability"
)

func withCorrelationIdMiddleware(ctx context.Context) (context.Context, error) {
	ctxWithCorrelationID := observability.ContextWithCorrelationID(ctx, observability.NewCorrelationID())
	return ctxWithCorrelationID, nil
}

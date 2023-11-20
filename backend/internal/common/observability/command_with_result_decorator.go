package observability

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/flowck/doberman/internal/common/logs"
)

type commandWithResultDecorator[C, R any] struct {
	base CommandWithResultHandler[C, R]
}

func NewCommandWithResultDecorator[C, R any](base CommandWithResultHandler[C, R], logger *logs.Logger, tracer trace.Tracer) commandWithResultDecorator[C, R] {
	return commandWithResultDecorator[C, R]{
		base: commandWithResultMetricsDecorator[C, R]{
			base: commandWithResultTracingDecorator[C, R]{
				base: commandWithResultLoggingDecorator[C, R]{
					base:   base,
					logger: logger,
				},
				// tracer: tracer,
			},
		},
	}
}

func (q commandWithResultDecorator[C, R]) Execute(ctx context.Context, cmd C) (result R, err error) {
	return q.base.Execute(ctx, cmd)
}

package observability

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/flowck/doberman/internal/common/logs"
)

type commandDecorator[C any] struct {
	base CommandHandler[C]
}

func NewCommandDecorator[C any](base CommandHandler[C], logger *logs.Logger, tracer trace.Tracer) commandDecorator[C] {
	return commandDecorator[C]{
		base: commandMetricsDecorator[C]{
			base: commandTracingDecorator[C]{
				base: commandLoggingDecorator[C]{
					base:   base,
					logger: logger,
				},
				// tracer: tracer,
			},
		},
	}
}

func (q commandDecorator[C]) Execute(ctx context.Context, cmd C) error {
	return q.base.Execute(ctx, cmd)
}

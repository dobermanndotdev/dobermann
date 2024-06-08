package cron

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/app/command"
)

func (s handlers) BulkCheckEndpoints(ctx context.Context) error {
	return s.application.Commands.BulkCheckEndpoints.Execute(ctx, command.BulkCheckEndpoints{})
}

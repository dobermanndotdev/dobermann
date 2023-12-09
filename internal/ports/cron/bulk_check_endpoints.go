package cron

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/app/command"
)

func (s handlers) BulkCheckEndpoints(ctx context.Context) error {
	return s.application.Commands.BulkCheckEndpoints.Execute(ctx, command.BulkCheckEndpoints{})
}

package events

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/app/command"
)

type Publisher struct{}

func NewPublisher() Publisher {
	return Publisher{}
}

func (p Publisher) PublishMonitorCreated(ctx context.Context, event command.MonitorCreatedEvent) error {
	return nil
}

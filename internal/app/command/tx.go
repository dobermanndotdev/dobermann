package command

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/domain/account"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
)

type TransactableAdapters struct {
	// IAM
	AccountRepository account.Repository
	UserRepository    account.UserRepository

	// Monitor
	MonitorRepository  monitor.Repository
	IncidentRepository monitor.IncidentRepository

	// Event publisher
	EventPublisher EventPublisher
}

type EventPublisher interface {
	PublishMonitorCreated(ctx context.Context, event MonitorCreatedEvent) error
	PublishIncidentCreated(ctx context.Context, event IncidentCreatedEvent) error
	PublishIncidentResolved(ctx context.Context, event IncidentResolvedEvent) error
	PublishEndpointCheckFailed(ctx context.Context, event EndpointCheckFailedEvent) error
	PublishEndpointCheckSucceeded(ctx context.Context, event EndpointCheckSucceededEvent) error
}

type TransactionProvider interface {
	Transact(ctx context.Context, f TransactFunc) error
}

type TransactFunc func(adapters TransactableAdapters) error

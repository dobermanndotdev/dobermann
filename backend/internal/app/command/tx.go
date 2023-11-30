package command

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
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
	PublishEndpointCheckFailed(ctx context.Context, event EndpointCheckFailedEvent) error
}

type TransactionProvider interface {
	Transact(ctx context.Context, f TransactFuncc) error
}

type TransactFuncc func(adapters TransactableAdapters) error

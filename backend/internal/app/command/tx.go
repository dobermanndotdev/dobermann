package command

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type TransactableAdapters struct {
	AccountRepository account.Repository
	MonitorRepository monitor.Repository
	UserRepository    account.UserRepository
	EventPublisher    EventPublisher
}

type EventPublisher interface {
	PublishMonitorCreated(ctx context.Context, event MonitorCreatedEvent) error
}

type TransactionProvider interface {
	Transact(ctx context.Context, f TransactFunc) error
}

type TransactFunc func(adapters TransactableAdapters) error

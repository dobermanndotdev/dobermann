package command

import (
	"github.com/flowck/doberman/internal/domain/account"
	"github.com/flowck/doberman/internal/domain/monitor"
)

type TransactableAdapters struct {
	EventPublisher    EventPublisher
	MonitorRepository monitor.Repository
	AccountRepository account.Repository
}

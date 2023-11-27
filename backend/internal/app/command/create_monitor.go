package command

import (
	"context"
	"fmt"

	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type CreateMonitor struct {
	Monitor *monitor.Monitor
}

type CreateMonitorHandler struct {
	txProvider TransactionProvider
}

func NewCreateMonitorHandler(txProvider TransactionProvider) CreateMonitorHandler {
	return CreateMonitorHandler{
		txProvider: txProvider,
	}
}

func (h CreateMonitorHandler) Execute(ctx context.Context, cmd CreateMonitor) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		err := adapters.MonitorRepository.Insert(ctx, cmd.Monitor)
		if err != nil {
			return err
		}

		err = adapters.EventPublisher.PublishMonitorCreated(ctx, MonitorCreatedEvent{
			ID:        cmd.Monitor.ID().String(),
			CreatedAt: cmd.Monitor.CreatedAt(),
		})
		if err != nil {
			return fmt.Errorf("unable to publish event: %v", err)
		}

		return nil
	})
}

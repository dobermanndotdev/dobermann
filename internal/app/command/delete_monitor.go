package command

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain"
)

type DeleteMonitor struct {
	ID domain.ID
}

type DeleteMonitorHandler struct {
	txProvider TransactionProvider
}

func NewDeleteMonitorHandler(txProvider TransactionProvider) DeleteMonitorHandler {
	return DeleteMonitorHandler{
		txProvider: txProvider,
	}
}

func (h DeleteMonitorHandler) Execute(ctx context.Context, cmd DeleteMonitor) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		return adapters.MonitorRepository.Delete(ctx, cmd.ID)
	})
}

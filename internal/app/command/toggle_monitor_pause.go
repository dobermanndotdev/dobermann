package command

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type ToggleMonitorPause struct {
	MonitorID domain.ID
	Pause     bool
}

type ToggleMonitorPauseHandler struct {
	txProvider TransactionProvider
}

func NewToggleMonitorPauseHandler(txProvider TransactionProvider) ToggleMonitorPauseHandler {
	return ToggleMonitorPauseHandler{
		txProvider: txProvider,
	}
}

func (h ToggleMonitorPauseHandler) Execute(ctx context.Context, cmd ToggleMonitorPause) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		return adapters.MonitorRepository.Update(ctx, cmd.MonitorID, func(foundMonitor *monitor.Monitor) error {
			if cmd.Pause {
				foundMonitor.Pause()
			} else {
				foundMonitor.UnPause()
			}

			return nil
		})
	})
}

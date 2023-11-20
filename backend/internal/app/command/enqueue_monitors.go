package command

import (
	"context"
	"time"

	"github.com/flowck/doberman/internal/common/tx"
	"github.com/flowck/doberman/internal/domain/monitor"
)

type EnqueueMonitors struct{}

type EnqueueMonitorsHandler struct {
	txProvider tx.TransactionProvider[TransactableAdapters]
}

func NewEnqueueMonitorsHandler(txProvider tx.TransactionProvider[TransactableAdapters]) EnqueueMonitorsHandler {
	return EnqueueMonitorsHandler{
		txProvider: txProvider,
	}
}

func (h EnqueueMonitorsHandler) Execute(ctx context.Context, cmd EnqueueMonitors) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		var monitorIDs []string
		err := adapters.MonitorRepository.UpdateAllForCheck(ctx, func(mos []*monitor.Monitor) error {
			for _, mo := range mos {
				mo.Enqueue()

				monitorIDs = append(monitorIDs, mo.ID().String())
			}

			return nil
		})

		if err != nil {
			return err
		}

		err = adapters.EventPublisher.PublishMonitorsEnqueued(ctx, MonitorsEnqueuedEvent{
			IDs:        monitorIDs,
			EnqueuedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		return nil
	})
}

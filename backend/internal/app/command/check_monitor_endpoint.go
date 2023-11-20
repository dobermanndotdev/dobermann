package command

import (
	"context"
	"fmt"
	"time"

	"github.com/flowck/doberman/internal/common/tx"
	"github.com/flowck/doberman/internal/domain"
	"github.com/flowck/doberman/internal/domain/monitor"
)

type CheckMonitorEndpoint struct {
	MonitorID domain.ID
}

type CheckResult struct {
	ResponseDuration time.Duration
	ResponseStatus   int
}

type CheckService interface {
	Check(ctx context.Context, endpoint monitor.Endpoint) (CheckResult, error)
}

type checkMonitorEndpointHandler struct {
	checkService CheckService
	txProvider   tx.TransactionProvider[TransactableAdapters]
}

func NewCheckMonitorEndpointHandler(txProvider tx.TransactionProvider[TransactableAdapters], checkService CheckService) checkMonitorEndpointHandler {
	return checkMonitorEndpointHandler{
		txProvider:   txProvider,
		checkService: checkService,
	}
}

func (h checkMonitorEndpointHandler) Execute(ctx context.Context, cmd CheckMonitorEndpoint) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		err := adapters.MonitorRepository.Update(ctx, cmd.MonitorID, func(mon *monitor.Monitor) error {
			checkResult, err := h.checkService.Check(ctx, mon.Endpoint())
			if err != nil {
				mon.SetStatusAsDown()

				if err = adapters.EventPublisher.PublishMonitorEndpointCheckFailed(ctx, mon.ID()); err != nil {
					return fmt.Errorf("unable to publish event: %v", err)
				}

				return nil // ack
			}

			fmt.Println("Response time", checkResult)
			err = adapters.EventPublisher.PublishMonitorEndpointCheckSucceeded(ctx, mon.ID())
			if err != nil {
				return fmt.Errorf("unable to publish event: %v", err)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("unable to update the monitor: %v", err)
		}

		return nil
	})
}

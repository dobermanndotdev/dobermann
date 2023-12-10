package command

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type BulkCheckEndpoints struct {
	FromRegion string
}

type BulkCheckEndpointsHandler struct {
	httpChecker httpChecker
	txProvider  TransactionProvider
}

func NewBulkCheckEndpointsHandler(httpChecker httpChecker, txProvider TransactionProvider) BulkCheckEndpointsHandler {
	return BulkCheckEndpointsHandler{
		txProvider:  txProvider,
		httpChecker: httpChecker,
	}
}

func (c BulkCheckEndpointsHandler) Execute(ctx context.Context, cmd BulkCheckEndpoints) error {
	return c.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		err := adapters.MonitorRepository.UpdateForCheck(ctx, func(foundMonitors []*monitor.Monitor) error {
			for _, foundMonitor := range foundMonitors {
				err := c.httpChecker.Check(ctx, foundMonitor.EndpointUrl())
				if err != nil {
					foundMonitor.SetEndpointCheckResult(false)
					err = adapters.EventPublisher.PublishEndpointCheckFailed(ctx, EndpointCheckFailedEvent{
						MonitorID: foundMonitor.ID().String(),
						At:        *foundMonitor.LastCheckedAt(),
					})
					if err != nil {
						return err
					}

					continue
				}

				foundMonitor.SetEndpointCheckResult(true)

				if !foundMonitor.HasIncidentUnresolved() {
					return nil
				}

				err = adapters.EventPublisher.PublishEndpointCheckSucceeded(ctx, EndpointCheckSucceededEvent{
					MonitorID: foundMonitor.ID().String(),
					At:        *foundMonitor.LastCheckedAt(),
				})
				if err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})
}

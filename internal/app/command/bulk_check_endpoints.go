package command

import (
	"context"
	"fmt"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type BulkCheckEndpoints struct {
}

type BulkCheckEndpointsHandler struct {
	httpChecker       httpChecker
	txProvider        TransactionProvider
	monitorRepository monitor.Repository
}

func NewBulkCheckEndpointsHandler(
	httpChecker httpChecker,
	txProvider TransactionProvider,
	monitorRepository monitor.Repository,
) BulkCheckEndpointsHandler {
	return BulkCheckEndpointsHandler{
		txProvider:        txProvider,
		httpChecker:       httpChecker,
		monitorRepository: monitorRepository,
	}
}

func (c BulkCheckEndpointsHandler) Execute(ctx context.Context, cmd BulkCheckEndpoints) error {
	checkResults := make(map[domain.ID]*monitor.CheckResult)

	err := c.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		return adapters.MonitorRepository.UpdateForCheck(ctx, func(foundMonitors []*monitor.Monitor) error {
			for _, foundMonitor := range foundMonitors {
				if foundMonitor.IsPaused() {
					continue
				}

				checkResult, err := c.httpChecker.Check(ctx, foundMonitor.EndpointUrl())
				if err != nil {
					return fmt.Errorf("error checking endpoint %s due to: %v", foundMonitor.EndpointUrl(), err)
				}

				checkResults[foundMonitor.ID()] = checkResult
				if checkResult.IsEndpointDown() {
					foundMonitor.MarkEndpointAsUp()
					err = adapters.EventPublisher.PublishEndpointCheckFailed(ctx, EndpointCheckFailedEvent{
						MonitorID: foundMonitor.ID().String(),
						At:        *foundMonitor.LastCheckedAt(),
					})
					if err != nil {
						return err
					}

					continue
				}

				foundMonitor.MarkEndpointAsUp()
				if foundMonitor.HasIncidentUnresolved() {
					err = adapters.EventPublisher.PublishEndpointCheckSucceeded(ctx, EndpointCheckSucceededEvent{
						MonitorID: foundMonitor.ID().String(),
						At:        *foundMonitor.LastCheckedAt(),
					})
					if err != nil {
						return err
					}
				}

				return nil
			}

			return nil
		})
	})

	if err != nil {
		return err
	}

	for checkedMonitorID, checkResult := range checkResults {
		err = c.monitorRepository.SaveCheckResult(ctx, checkedMonitorID, checkResult)
		if err != nil {
			return fmt.Errorf("unable to save the check result of monitor with id %s: %v", checkedMonitorID, err)
		}
	}

	return nil
}

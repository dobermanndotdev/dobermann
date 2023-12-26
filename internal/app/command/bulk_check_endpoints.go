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

	err := c.txProvider.Transact(ctx, c.checkMatchedEndpoints(ctx, &checkResults))
	if err != nil {
		return err
	}

	err = c.saveCheckResults(ctx, &checkResults)
	if err != nil {
		return err
	}

	return nil
}

func (c BulkCheckEndpointsHandler) saveCheckResults(ctx context.Context, results *map[domain.ID]*monitor.CheckResult) error {
	var err error
	for checkedMonitorID, checkResult := range *results {
		err = c.monitorRepository.SaveCheckResult(ctx, checkedMonitorID, checkResult)
		if err != nil {
			return fmt.Errorf("unable to save the check result of monitor with id %s: %v", checkedMonitorID, err)
		}
	}

	return nil
}

func (c BulkCheckEndpointsHandler) checkMatchedEndpoints(
	ctx context.Context,
	results *map[domain.ID]*monitor.CheckResult,
) func(adapters TransactableAdapters) error {
	return func(adapters TransactableAdapters) error {
		return adapters.MonitorRepository.UpdateForCheck(ctx, func(foundMonitors []*monitor.Monitor) error {
			for _, foundMonitor := range foundMonitors {
				if foundMonitor.IsPaused() {
					continue
				}

				checkResult, err := c.httpChecker.Check(ctx, foundMonitor.EndpointUrl())
				if err != nil {
					return fmt.Errorf("error checking endpoint %s due to: %v", foundMonitor.EndpointUrl(), err)
				}

				(*results)[foundMonitor.ID()] = checkResult

				if checkResult.IsEndpointDown() {
					err = c.handleEndpointDown(ctx, adapters, foundMonitor)
					if err != nil {
						return err
					}

					continue
				}

				err = c.handleEndpointIsUp(ctx, adapters, foundMonitor)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}
}

func (c BulkCheckEndpointsHandler) handleEndpointDown(
	ctx context.Context,
	adapters TransactableAdapters,
	m *monitor.Monitor,
) error {
	m.MarkEndpointAsDown()

	return adapters.EventPublisher.PublishEndpointCheckFailed(ctx, EndpointCheckFailedEvent{
		MonitorID: m.ID().String(),
		At:        *m.LastCheckedAt(),
	})
}

func (c BulkCheckEndpointsHandler) handleEndpointIsUp(
	ctx context.Context,
	adapters TransactableAdapters,
	m *monitor.Monitor,
) error {
	m.MarkEndpointAsUp()

	if !m.HasIncidentUnresolved() {
		return nil
	}

	return adapters.EventPublisher.PublishEndpointCheckSucceeded(ctx, EndpointCheckSucceededEvent{
		MonitorID: m.ID().String(),
		At:        *m.LastCheckedAt(),
	})
}

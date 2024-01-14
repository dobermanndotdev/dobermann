package command

import (
	"context"
	"fmt"

	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type BulkCheckEndpoints struct {
}

type BulkCheckEndpointsHandler struct {
	httpChecker       httpChecker
	txProvider        TransactionProvider
	eventPublisher    EventPublisher
	monitorRepository monitor.Repository
}

func NewBulkCheckEndpointsHandler(
	httpChecker httpChecker,
	txProvider TransactionProvider,
	eventPublisher EventPublisher,
	monitorRepository monitor.Repository,
) BulkCheckEndpointsHandler {
	return BulkCheckEndpointsHandler{
		txProvider:        txProvider,
		httpChecker:       httpChecker,
		eventPublisher:    eventPublisher,
		monitorRepository: monitorRepository,
	}
}

type result struct {
	Monitor     *monitor.Monitor
	CheckResult *monitor.CheckResult
}

func (c BulkCheckEndpointsHandler) Execute(ctx context.Context, cmd BulkCheckEndpoints) error {
	return c.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		return adapters.MonitorRepository.UpdateForCheck(ctx, func(foundMonitors []*monitor.Monitor) error {
			for _, foundMonitor := range foundMonitors {
				if foundMonitor.IsPaused() {
					continue
				}

				checkResult, err := c.httpChecker.Check(ctx, foundMonitor.EndpointUrl())
				if err != nil {
					return fmt.Errorf("error checking endpoint %s due to: %v", foundMonitor.EndpointUrl(), err)
				}

				err = c.monitorRepository.SaveCheckResult(ctx, foundMonitor.ID(), checkResult)
				if err != nil {
					return fmt.Errorf("unable to save the check result: %v", err)
				}

				if checkResult.IsEndpointDown() {
					err = c.handleEndpointDown(ctx, foundMonitor, checkResult)
					if err != nil {
						return err
					}

					continue
				}

				err = c.handleEndpointIsUp(ctx, foundMonitor)
				if err != nil {
					return err
				}
			}

			return nil
		})
	})
}

func (c BulkCheckEndpointsHandler) handleEndpointDown(
	ctx context.Context,
	m *monitor.Monitor,
	checkResult *monitor.CheckResult,
) error {
	m.MarkEndpointAsDown()

	return c.eventPublisher.PublishEndpointCheckFailed(ctx, EndpointCheckFailedEvent{
		MonitorID:     m.ID().String(),
		CheckedURL:    m.EndpointUrl(),
		At:            *m.LastCheckedAt(),
		CheckResultID: checkResult.ID().String(),
	})
}

func (c BulkCheckEndpointsHandler) handleEndpointIsUp(ctx context.Context, m *monitor.Monitor) error {
	m.MarkEndpointAsUp()

	if !m.HasIncidentUnresolved() {
		return nil
	}

	return c.eventPublisher.PublishEndpointCheckSucceeded(ctx, EndpointCheckSucceededEvent{
		MonitorID: m.ID().String(),
		At:        *m.LastCheckedAt(),
	})
}

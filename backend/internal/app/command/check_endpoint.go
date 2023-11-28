package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type CheckEndpoint struct {
	MonitorID domain.ID
}

type CheckEndpointHandler struct {
	httpChecker       httpChecker
	eventPublisher    EventPublisher
	monitorRepository monitor.Repository
}

type httpChecker interface {
	Check(ctx context.Context, endpointUrl string) error
}

func NewCheckEndpointHandler(
	httpChecker httpChecker,
	monitorRepository monitor.Repository,
	eventPublisher EventPublisher,
) CheckEndpointHandler {
	return CheckEndpointHandler{
		httpChecker:       httpChecker,
		eventPublisher:    eventPublisher,
		monitorRepository: monitorRepository,
	}
}

func (c CheckEndpointHandler) Execute(ctx context.Context, cmd CheckEndpoint) error {
	checkSucceeded := false

	err := c.monitorRepository.Update(ctx, cmd.MonitorID, func(m *monitor.Monitor) error {
		err := c.httpChecker.Check(ctx, m.EndpointUrl())
		if errors.Is(err, monitor.ErrEndpointIsDown) {
			m.SetEndpointCheckResult(false)
			return nil
		}

		if err != nil {
			return fmt.Errorf("error checking endpoint %s due to: %v", m.EndpointUrl(), err)
		}

		checkSucceeded = true
		m.SetEndpointCheckResult(true)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error updating monitor during check: %v", err)
	}

	if !checkSucceeded {
		err = c.eventPublisher.PublishEndpointCheckFailed(ctx, EndpointCheckFailed{
			At:        time.Now(),
			MonitorID: cmd.MonitorID.String(),
		})

		if err != nil {
			return fmt.Errorf("error publishing event EndpointCheckFailed: %v", err)
		}
	}

	return nil
}

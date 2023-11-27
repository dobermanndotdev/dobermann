package command

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type CheckEndpoint struct {
	MonitorID domain.ID
}

type CheckEndpointHandler struct {
	httpChecker       httpChecker
	monitorRepository monitor.Repository
}

type httpChecker interface {
	Check(ctx context.Context, endpointUrl string) error
}

func NewCheckEndpointHandler(httpChecker httpChecker, monitorRepository monitor.Repository) CheckEndpointHandler {
	return CheckEndpointHandler{
		httpChecker:       httpChecker,
		monitorRepository: monitorRepository,
	}
}

func (c CheckEndpointHandler) Execute(ctx context.Context, cmd CheckEndpoint) error {
	err := c.monitorRepository.Update(ctx, cmd.MonitorID, func(m *monitor.Monitor) error {
		err := c.httpChecker.Check(ctx, m.EndpointUrl())
		if err != nil {
			m.SetEndpointCheckResult(false)
			return nil
		}

		m.SetEndpointCheckResult(true)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

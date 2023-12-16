package command

import (
	"context"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type EditMonitor struct {
	ID                     domain.ID
	EndpointUrl            string
	CheckIntervalInSeconds time.Duration
}

type EditMonitorHandler struct {
	monitorRepository monitor.Repository
}

func NewEditMonitorHandler(monitorRepository monitor.Repository) EditMonitorHandler {
	return EditMonitorHandler{
		monitorRepository: monitorRepository,
	}
}

func (h EditMonitorHandler) Execute(ctx context.Context, cmd EditMonitor) error {
	return h.monitorRepository.Update(ctx, cmd.ID, func(foundMonitor *monitor.Monitor) error {
		return foundMonitor.Edit(cmd.EndpointUrl, cmd.CheckIntervalInSeconds)
	})
}

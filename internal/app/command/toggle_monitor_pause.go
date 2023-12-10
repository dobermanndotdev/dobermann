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
	monitorRepository monitor.Repository
}

func NewToggleMonitorPauseHandler(monitorRepository monitor.Repository) ToggleMonitorPauseHandler {
	return ToggleMonitorPauseHandler{
		monitorRepository: monitorRepository,
	}
}

func (h ToggleMonitorPauseHandler) Execute(ctx context.Context, cmd ToggleMonitorPause) error {
	return h.monitorRepository.Update(ctx, cmd.MonitorID, func(foundMonitor *monitor.Monitor) error {
		if cmd.Pause {
			foundMonitor.Pause()
		} else {
			foundMonitor.UnPause()
		}

		return nil
	})
}

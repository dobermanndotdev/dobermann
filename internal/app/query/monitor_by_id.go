package query

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
)

type MonitorByID struct {
	ID domain.ID
}

type MonitorByIdHandler struct {
	monitorRepository monitor.Repository
}

func NewMonitorByIdHandler(monitorRepository monitor.Repository) MonitorByIdHandler {
	return MonitorByIdHandler{
		monitorRepository: monitorRepository,
	}
}

func (h MonitorByIdHandler) Execute(ctx context.Context, q MonitorByID) (*monitor.Monitor, error) {
	return h.monitorRepository.FindByID(ctx, q.ID)
}

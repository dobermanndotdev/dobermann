package query

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type MonitorByID struct {
	AccountID domain.ID
	ID        domain.ID
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
	return h.monitorRepository.FindByID(ctx, q.AccountID, q.ID)
}

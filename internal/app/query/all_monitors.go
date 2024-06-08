package query

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
)

type AllMonitors struct {
	AccountID domain.ID
	Params    PaginationParams
}

type monitorsFinder interface {
	FindAll(ctx context.Context, accID domain.ID, params PaginationParams) (PaginatedResult[*monitor.Monitor], error)
}

type AllMonitorsHandler struct {
	monitorsFinder monitorsFinder
}

func NewAllMonitorsHandler(monitorsFinder monitorsFinder) AllMonitorsHandler {
	return AllMonitorsHandler{
		monitorsFinder: monitorsFinder,
	}
}

func (h AllMonitorsHandler) Execute(ctx context.Context, q AllMonitors) (PaginatedResult[*monitor.Monitor], error) {
	return h.monitorsFinder.FindAll(ctx, q.AccountID, q.Params)
}

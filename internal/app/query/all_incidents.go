package query

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
)

type AllIncidents struct {
	AccountID domain.ID
	Params    PaginationParams
}

type incidentsFinder interface {
	FindAll(ctx context.Context, accID domain.ID, params PaginationParams) (PaginatedResult[*monitor.Incident], error)
}

type AllIncidentsHandler struct {
	incidentsFinder incidentsFinder
}

func NewAllIncidentsHandler(incidentsFinder incidentsFinder) AllIncidentsHandler {
	return AllIncidentsHandler{
		incidentsFinder: incidentsFinder,
	}
}

func (h AllIncidentsHandler) Execute(ctx context.Context, q AllIncidents) (PaginatedResult[*monitor.Incident], error) {
	return h.incidentsFinder.FindAll(ctx, q.AccountID, q.Params)
}

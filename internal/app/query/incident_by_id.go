package query

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
)

type IncidentByID struct {
	ID domain.ID
}

type IncidentByIdHandler struct {
	incidentRepository monitor.IncidentRepository
}

func NewIncidentByIdHandler(incidentRepository monitor.IncidentRepository) IncidentByIdHandler {
	return IncidentByIdHandler{
		incidentRepository: incidentRepository,
	}
}

func (h IncidentByIdHandler) Execute(ctx context.Context, q IncidentByID) (*monitor.Incident, error) {
	return h.incidentRepository.FindByID(ctx, q.ID)
}

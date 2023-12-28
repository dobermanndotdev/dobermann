package query

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
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

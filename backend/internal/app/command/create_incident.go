package command

import (
	"context"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type CreateIncident struct {
	MonitorID domain.ID
}

type CreateIncidentHandler struct {
	incidentRepository monitor.IncidentRepository
}

func NewCreateIncidentHandler(incidentRepository monitor.IncidentRepository) CreateIncidentHandler {
	return CreateIncidentHandler{
		incidentRepository: incidentRepository,
	}
}

func (h CreateIncidentHandler) Execute(ctx context.Context, cmd CreateIncident) error {
	incident, err := monitor.NewIncident(domain.NewID(), time.Now().UTC(), nil)
	if err != nil {
		return err
	}

	return h.incidentRepository.Create(ctx, cmd.MonitorID, incident)
}

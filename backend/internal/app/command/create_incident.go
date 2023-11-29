package command

import (
	"context"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type CreateIncident struct {
	MonitorID domain.ID
}

type CreateIncidentHandler struct {
	monitorRepository  monitor.Repository
	userRepository     account.UserRepository
	subscriberNotifier subscriberNotifier
	incidentRepository monitor.IncidentRepository
}

type subscriberNotifier interface {
	SendEmailAboutIncident(ctx context.Context, user *account.User, incident *monitor.Incident) error
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

	err = h.incidentRepository.Create(ctx, cmd.MonitorID, incident)
	if err != nil {
		return err
	}

	// m, err := h.monitorRepository.

	// h.userRepository.FindByID(ctx, c)

	return nil
}

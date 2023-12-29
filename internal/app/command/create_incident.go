package command

import (
	"context"
	"fmt"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type CreateIncident struct {
	MonitorID  domain.ID
	CheckedURL string
	Details    monitor.IncidentDetails
}

type CreateIncidentHandler struct {
	txProvider TransactionProvider
}

func NewCreateIncidentHandler(txProvider TransactionProvider) CreateIncidentHandler {
	return CreateIncidentHandler{
		txProvider: txProvider,
	}
}

// TODO: Test this command's logic ;)
func (h CreateIncidentHandler) Execute(ctx context.Context, cmd CreateIncident) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		foundMonitor, err := adapters.MonitorRepository.FindByID(ctx, cmd.MonitorID)
		if err != nil {
			return err
		}

		if foundMonitor.HasIncidentUnresolved() {
			return nil
		}

		incident, err := monitor.NewIncident(domain.NewID(), foundMonitor.ID(), nil, time.Now().UTC(), cmd.CheckedURL, nil, cmd.Details)
		if err != nil {
			return err
		}

		err = adapters.IncidentRepository.Create(ctx, incident)
		if err != nil {
			return fmt.Errorf("unable to save incident: %v", err)
		}

		err = adapters.EventPublisher.PublishIncidentCreated(ctx, IncidentCreatedEvent{
			MonitorID:  cmd.MonitorID.String(),
			IncidentID: incident.ID().String(),
			At:         incident.CreatedAt(),
		})
		if err != nil {
			return fmt.Errorf("unable to publish event IncidentCreatedEvent: %v", err)
		}

		return nil
	})
}

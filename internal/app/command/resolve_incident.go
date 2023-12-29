package command

import (
	"context"
	"fmt"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type ResolveIncident struct {
	MonitorID domain.ID
}

type ResolveIncidentHandler struct {
	txProvider TransactionProvider
}

func NewResolveIncidentHandler(txProvider TransactionProvider) ResolveIncidentHandler {
	return ResolveIncidentHandler{
		txProvider: txProvider,
	}
}

func (h ResolveIncidentHandler) Execute(ctx context.Context, cmd ResolveIncident) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		foundMonitor, err := adapters.MonitorRepository.FindByID(ctx, cmd.MonitorID)
		if err != nil {
			return err
		}

		incident := foundMonitor.IncidentUnresolved()
		if incident == nil {
			return nil
		}

		var incidentAction *monitor.IncidentAction
		incidentAction, err = monitor.NewIncidentAction(domain.NewID(), nil, time.Now(), "", monitor.IncidentActionTypeResolved)
		if err != nil {
			return err
		}

		err = adapters.IncidentRepository.AppendIncidentAction(ctx, incident.ID(), incidentAction)
		if err != nil {
			return fmt.Errorf("unable to append incident action: %v", err)
		}

		err = adapters.IncidentRepository.Update(ctx, incident.ID(), func(incident *monitor.Incident) error {
			incident.Resolve()

			return nil
		})

		if err != nil {
			return err
		}

		err = adapters.EventPublisher.PublishIncidentResolved(ctx, IncidentResolvedEvent{
			MonitorID:  cmd.MonitorID.String(),
			IncidentID: incident.ID().String(),
			At:         time.Now().UTC(),
		})
		if err != nil {
			return err
		}

		return nil
	})
}

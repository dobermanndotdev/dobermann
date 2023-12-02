package command

import (
	"context"
	"fmt"
	"time"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type ResolveIncidents struct {
	MonitorID domain.ID
}

type ResolveIncidentsHandler struct {
	txProvider TransactionProvider
}

func NewResolveIncidentsHandler(txProvider TransactionProvider) ResolveIncidentsHandler {
	return ResolveIncidentsHandler{
		txProvider: txProvider,
	}
}

func (h ResolveIncidentsHandler) Execute(ctx context.Context, cmd ResolveIncidents) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		foundMonitor, err := adapters.MonitorRepository.FindByID(ctx, cmd.MonitorID)
		if err != nil {
			return err
		}

		for _, incident := range foundMonitor.Incidents() {
			if incident.IsResolved() {
				continue
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

			err = adapters.IncidentRepository.Update(ctx, incident.ID(), cmd.MonitorID, func(incident *monitor.Incident) error {
				incident.Resolve()

				return nil
			})

			if err != nil {
				return err
			}
		}

		err = adapters.EventPublisher.PublishIncidentResolved(ctx, IncidentResolvedEvent{
			MonitorID: cmd.MonitorID.String(),
			At:        time.Now().UTC(),
		})
		if err != nil {
			return err
		}

		return nil
	})
}

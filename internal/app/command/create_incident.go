package command

import (
	"context"
	"fmt"
	"time"

	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
)

type CreateIncident struct {
	MonitorID     domain.ID
	CheckResultID domain.ID
	CheckedUrl    string
}

type checkResultFinder interface {
	FindById(ctx context.Context, id domain.ID) (*monitor.CheckResult, error)
}

type CreateIncidentHandler struct {
	txProvider        TransactionProvider
	checkResultFinder checkResultFinder
}

func NewCreateIncidentHandler(txProvider TransactionProvider, checkResultFinder checkResultFinder) CreateIncidentHandler {
	return CreateIncidentHandler{
		txProvider:        txProvider,
		checkResultFinder: checkResultFinder,
	}
}

// TODO: Test this command's logic ;)
func (h CreateIncidentHandler) Execute(ctx context.Context, cmd CreateIncident) error {
	checkResult, err := h.checkResultFinder.FindById(ctx, cmd.CheckResultID)
	if err != nil {
		return err
	}

	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		foundMonitor, err := adapters.MonitorRepository.FindByID(ctx, cmd.MonitorID)
		if err != nil {
			return err
		}

		if foundMonitor.HasIncidentUnresolved() {
			return nil
		}

		incident, err := monitor.NewIncident(
			domain.NewID(),
			foundMonitor.ID(),
			nil,
			time.Now().UTC(),
			cmd.CheckedUrl,
			nil,
			fmt.Sprintf("Monitor with url %s is unresponsive", cmd.CheckedUrl),
			checkResult.StatusCode(),
		)
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

package command

import (
	"context"
	"fmt"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type NotifyMonitorSubscribersOnIncident struct {
	IncidentID domain.ID
	MonitorID  domain.ID
}

type NotifyMonitorSubscribersOnIncidentHandler struct {
	txProvider         TransactionProvider
	subscriberNotifier subscriberNotifier
}

type subscriberNotifier interface {
	SendEmailAboutIncident(context.Context, *account.User, *monitor.Monitor, *monitor.Incident) error
}

func NewNotifyMonitorSubscribersOnIncidentHandler(
	txProvider TransactionProvider,
	subscriberNotifier subscriberNotifier,
) NotifyMonitorSubscribersOnIncidentHandler {
	return NotifyMonitorSubscribersOnIncidentHandler{
		txProvider:         txProvider,
		subscriberNotifier: subscriberNotifier,
	}
}

func (h NotifyMonitorSubscribersOnIncidentHandler) Execute(
	ctx context.Context,
	cmd NotifyMonitorSubscribersOnIncident,
) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		m, err := adapters.MonitorRepository.FindByID(ctx, cmd.MonitorID)
		if err != nil {
			return fmt.Errorf("error while trying to find monitor via id %s :%v", cmd.MonitorID, err)
		}

		incident, err := adapters.IncidentRepository.FindByID(ctx, cmd.IncidentID)
		if err != nil {
			return fmt.Errorf("error while trying to find incident via id %s :%v", cmd.IncidentID, err)
		}

		var user *account.User
		for _, subscriber := range m.Subscribers() {
			user, err = adapters.UserRepository.FindByID(ctx, subscriber.UserID())
			if err != nil {
				return fmt.Errorf("error while trying to find user via id %s :%v", subscriber.UserID(), err)
			}

			//WARN: Error prone due to being an external service
			err = h.subscriberNotifier.SendEmailAboutIncident(ctx, user, m, incident)
			if err != nil {
				return fmt.Errorf("error while trying to notify user with id %s :%v", subscriber.UserID(), err)
			}
		}

		return nil
	})
}

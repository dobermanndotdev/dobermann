package command

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type NotifyOnIncidentResolved struct {
	IncidentID domain.ID
	MonitorID  domain.ID
}

type incidentResolvedNotifier interface {
	SendEmailIncidentResolution(ctx context.Context, user *account.User, m *monitor.Monitor, incidentID domain.ID) error
}

type NotifyOnIncidentResolvedHandler struct {
	monitorRepository  monitor.Repository
	userRepository     account.UserRepository
	subscriberNotifier incidentResolvedNotifier
}

func NewNotifyOnIncidentResolvedHandler(
	monitorRepository monitor.Repository,
	userRepository account.UserRepository,
	subscriberNotifier incidentResolvedNotifier,
) NotifyOnIncidentResolvedHandler {
	return NotifyOnIncidentResolvedHandler{
		userRepository:     userRepository,
		monitorRepository:  monitorRepository,
		subscriberNotifier: subscriberNotifier,
	}
}

func (h NotifyOnIncidentResolvedHandler) Execute(ctx context.Context, cmd NotifyOnIncidentResolved) error {
	foundMonitor, err := h.monitorRepository.FindByID(ctx, cmd.MonitorID)
	if err != nil {
		return err
	}

	var user *account.User
	for _, subscriber := range foundMonitor.Subscribers() {
		user, err = h.userRepository.FindByID(ctx, subscriber.UserID())
		if err != nil {
			return err
		}

		err = h.subscriberNotifier.SendEmailIncidentResolution(ctx, user, foundMonitor, cmd.IncidentID)
		if err != nil {
			return err
		}
	}

	return nil
}

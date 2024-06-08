package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/dobermanndotdev/dobermann/internal/adapters/events"
	"github.com/dobermanndotdev/dobermann/internal/app"
	"github.com/dobermanndotdev/dobermann/internal/app/command"
	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type IncidentCreatedHandler struct {
	application *app.App
}

func (h IncidentCreatedHandler) HandlerName() string {
	return "IncidentCreated_NotifyMonitorSubscribersOnIncident_Handler"
}

func (h IncidentCreatedHandler) EventName() string {
	return events.IncidentCreatedEvent{}.EventName()
}

func (h IncidentCreatedHandler) Handle(m *message.Message) error {
	event, err := events.NewEventFromMessage[events.IncidentCreatedEvent](m)
	if err != nil {
		return err
	}

	monitorID, err := domain.NewIdFromString(event.MonitorID)
	if err != nil {
		return err
	}

	incidentID, err := domain.NewIdFromString(event.IncidentID)
	if err != nil {
		return err
	}

	cmd := command.NotifyMonitorSubscribersOnIncident{
		IncidentID: incidentID,
		MonitorID:  monitorID,
	}
	err = h.application.Commands.NotifyMonitorSubscribersOnIncident.Execute(m.Context(), cmd)
	if err != nil {
		return err
	}

	return nil
}

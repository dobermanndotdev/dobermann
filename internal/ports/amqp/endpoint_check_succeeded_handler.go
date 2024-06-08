package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/dobermanndotdev/dobermann/internal/adapters/events"
	"github.com/dobermanndotdev/dobermann/internal/app"
	"github.com/dobermanndotdev/dobermann/internal/app/command"
	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type EndpointCheckSucceededHandler struct {
	application *app.App
}

func (e EndpointCheckSucceededHandler) HandlerName() string {
	return "EndpointCheckSucceeded_ResolveIncidentsHandler"
}

func (e EndpointCheckSucceededHandler) EventName() string {
	return events.EndpointCheckSucceededEvent{}.EventName()
}

func (e EndpointCheckSucceededHandler) Handle(m *message.Message) error {
	event, err := events.NewEventFromMessage[events.EndpointCheckSucceededEvent](m)
	if err != nil {
		return err
	}

	monitorID, err := domain.NewIdFromString(event.MonitorID)
	if err != nil {
		return err
	}

	err = e.application.Commands.ResolveIncident.Execute(m.Context(), command.ResolveIncident{
		MonitorID: monitorID,
	})
	if err != nil {
		return err
	}

	return nil
}

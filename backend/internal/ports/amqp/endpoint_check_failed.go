package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/domain"
)

type EndpointCheckFailedHandler struct {
	application *app.App
}

func (e EndpointCheckFailedHandler) HandlerName() string {
	return "EndpointCheckFailedHandler"
}

func (e EndpointCheckFailedHandler) EventName() string {
	return events.EndpointCheckFailed{}.EventName()
}

func (e EndpointCheckFailedHandler) Handle(m *message.Message) error {
	event, err := events.NewEventFromMessage[events.EndpointCheckFailed](m)
	if err != nil {
		return err
	}

	monitorID, err := domain.NewIdFromString(event.MonitorID)
	if err != nil {
		return err
	}

	err = e.application.Commands.CreateIncident.Execute(m.Context(), command.CreateIncident{
		MonitorID: monitorID,
	})
	if err != nil {
		return err
	}

	return nil
}

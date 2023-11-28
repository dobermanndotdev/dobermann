package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/domain"
)

type MonitorCreatedHandler struct {
	application *app.App
}

func (h MonitorCreatedHandler) HandlerName() string {
	return "MonitorCreatedHandler"
}

func (h MonitorCreatedHandler) EventName() string {
	return events.MonitorCreatedEvent{}.EventName()
}

func (h MonitorCreatedHandler) Handle(m *message.Message) error {
	event, err := events.NewMonitorCreatedEventFromMessage(m)
	if err != nil {
		return err
	}

	monitorID, err := domain.NewIdFromString(event.ID)
	if err != nil {
		return err
	}

	err = h.application.Commands.CheckEndpoint.Execute(m.Context(), command.CheckEndpoint{
		MonitorID: monitorID,
	})
	if err != nil {
		return err
	}

	return nil
}

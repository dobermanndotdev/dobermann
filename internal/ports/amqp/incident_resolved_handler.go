package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/app"
)

type IncidentResolvedHandler struct {
	application *app.App
}

func (e IncidentResolvedHandler) HandlerName() string {
	return "IncidentResolved_Handler"
}

func (e IncidentResolvedHandler) EventName() string {
	return events.IncidentResolvedEvent{}.EventName()
}

func (e IncidentResolvedHandler) Handle(m *message.Message) error {
	return nil
}

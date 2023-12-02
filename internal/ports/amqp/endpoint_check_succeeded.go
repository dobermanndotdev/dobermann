package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/app"
)

type EndpointCheckSucceededHandler struct {
	application *app.App
}

func (e EndpointCheckSucceededHandler) HandlerName() string {
	return "EndpointCheckSucceeded_Handler"
}

func (e EndpointCheckSucceededHandler) EventName() string {
	return events.EndpointCheckSucceededEvent{}.EventName()
}

func (e EndpointCheckSucceededHandler) Handle(m *message.Message) error {
	// Do nothing
	return nil
}

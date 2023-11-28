package amqp

import (
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/flowck/dobermann/backend/internal/app"
)

type Handler interface {
	HandlerName() string
	EventName() string
	Handle(m *message.Message) error
}

func NewHandlers(application *app.App) []Handler {
	return []Handler{
		MonitorCreatedHandler{application: application},
		EndpointCheckFailedHandler{application: application},
	}
}
